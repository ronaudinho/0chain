package storagesc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	chainState "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/state"
	"0chain.net/chaincore/tokenpool"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
	"0chain.net/core/datastore"
	"0chain.net/core/util"
)

// lock request

type lockRequest struct {
	Duration time.Duration `json:"duration"`
}

func (lr *lockRequest) decode(input []byte) error {
	return json.Unmarshal(input, lr)
}

// unlock request

type unlockRequest struct {
	PoolID datastore.Key `json:"pool_id"`
}

func (ur *unlockRequest) decode(input []byte) error {
	return json.Unmarshal(input, ur)
}

// read pool (a locked tokens for a duration)

type readPool struct {
	*tokenpool.ZcnLockingPool `json:"pool"`
}

func newReadPool() *readPool {
	return &readPool{ZcnLockingPool: &tokenpool.ZcnLockingPool{}}
}

func (rp *readPool) encode() (b []byte) {
	var err error
	if b, err = json.Marshal(rp); err != nil {
		panic(err) // must never happens
	}
	return
}

func (rp *readPool) decode(input []byte) (err error) {

	type readPoolJSON struct {
		Pool json.RawMessage `json:"pool"`
	}

	var readPoolVal readPoolJSON
	if err = json.Unmarshal(input, &readPoolVal); err != nil {
		return
	}

	if len(readPoolVal.Pool) == 0 {
		return // no data given
	}

	err = rp.ZcnLockingPool.Decode(readPoolVal.Pool, &tokenLock{})
	return
}

func (rp *readPool) stat(tp time.Time) (stat *readPoolStat, err error) {

	stat = new(readPoolStat)

	if err = stat.decode(rp.LockStats(tp)); err != nil {
		return nil, err
	}

	stat.ID = rp.ID
	stat.Locked = rp.IsLocked(tp)
	stat.Balance = rp.Balance

	return
}

// readPools -- set of locked tokens for a duration

// readPools of a user
type readPools struct {
	Pools map[datastore.Key]*readPool `json:"pools"`
}

func newReadPools() (rps *readPools) {
	rps = new(readPools)
	rps.Pools = make(map[datastore.Key]*readPool)
	return
}

func (rps *readPools) Encode() (b []byte) {
	var err error
	if b, err = json.Marshal(rps); err != nil {
		panic(err) // must never happens
	}
	return
}

func (rps *readPools) Decode(input []byte) (err error) {
	type readPoolsJSON struct {
		Pools map[string]json.RawMessage `json:"pools"`
	}
	var in readPoolsJSON
	if err = json.Unmarshal(input, &in); err != nil {
		return
	}
	for _, raw := range in.Pools {
		var tempPool = newReadPool()
		if err = tempPool.decode(raw); err != nil {
			return
		}
		rps.addPool(tempPool)
	}
	return
}

func readPoolsKey(scKey, clientID string) datastore.Key {
	return datastore.Key(scKey + ":readpool:" + clientID)
}

func (rps *readPools) addPool(rp *readPool) (err error) {
	if _, ok := rps.Pools[rp.ID]; ok {
		return errors.New("user already has this read pool")
	}
	rps.Pools[rp.ID] = rp
	return
}

func (rps *readPools) delPool(id datastore.Key) {
	delete(rps.Pools, id)
}

func (rps *readPools) save(sscKey, clientID string,
	balances chainState.StateContextI) (err error) {

	_, err = balances.InsertTrieNode(readPoolsKey(sscKey, clientID), rps)
	return
}

func (rps *readPools) moveToBlobber(now common.Timestamp, sp *stakePool,
	value state.Balance) (err error) {

	var tp = common.ToTime(now)

	for k, rp := range rps.Pools {
		if value == 0 {
			break
		}
		if !rp.IsLocked(tp) {
			continue
		}
		var move state.Balance
		if rp.Balance < value || rp.Balance == value {
			move = rp.Balance
			delete(rps.Pools, k)
		} else {
			move = value
		}
		if _, _, err = rp.TransferTo(sp.Unlocked, move, nil); err != nil {
			break
		}
		value -= move // decrease
	}

	if err != nil {
		return
	}

	if value != 0 {
		return errors.New("not enough tokens in read pool")
	}

	return
}

// stat

type readPoolStats struct {
	Stats []*readPoolStat `json:"stats"`
}

func (stats *readPoolStats) encode() (b []byte) {
	var err error
	if b, err = json.Marshal(stats); err != nil {
		panic(err) // must never happens
	}
	return
}

func (stats *readPoolStats) decode(input []byte) error {
	return json.Unmarshal(input, stats)
}

func (stats *readPoolStats) addStat(stat *readPoolStat) {
	stats.Stats = append(stats.Stats, stat)
}

type readPoolStat struct {
	ID        datastore.Key    `json:"pool_id"`
	StartTime common.Timestamp `json:"start_time"`
	Duartion  time.Duration    `json:"duration"`
	TimeLeft  time.Duration    `json:"time_left"`
	Locked    bool             `json:"locked"`
	Balance   state.Balance    `json:"balance"`
}

func (stat *readPoolStat) encode() (b []byte) {
	var err error
	if b, err = json.Marshal(stat); err != nil {
		panic(err) // must never happens
	}
	return
}

func (stat *readPoolStat) decode(input []byte) error {
	return json.Unmarshal(input, stat)
}

type tokenLock struct {
	StartTime common.Timestamp `json:"start_time"`
	Duration  time.Duration    `json:"duration"`
	Owner     datastore.Key    `json:"owner"`
}

func (tl tokenLock) IsLocked(entity interface{}) bool {
	if tm, ok := entity.(time.Time); ok {
		return tm.Sub(common.ToTime(tl.StartTime)) < tl.Duration
	}
	return true
}

func (tl tokenLock) LockStats(entity interface{}) []byte {
	if tm, ok := entity.(time.Time); ok {
		var stat readPoolStat
		stat.StartTime = tl.StartTime
		stat.Duartion = tl.Duration
		stat.TimeLeft = (tl.Duration - tm.Sub(common.ToTime(tl.StartTime)))
		stat.Locked = tl.IsLocked(tm)
		return stat.encode()
	}
	return nil
}

//
// smart contract methods
//

// getReadPoolsBytes of a client
func (ssc *StorageSmartContract) getReadPoolsBytes(clientID datastore.Key,
	balances chainState.StateContextI) (b []byte, err error) {

	var val util.Serializable
	val, err = balances.GetTrieNode(readPoolsKey(ssc.ID, clientID))
	if err != nil {
		return
	}
	return val.Encode(), nil
}

// getReadPools of current client
func (ssc *StorageSmartContract) getReadPools(clientID datastore.Key,
	balances chainState.StateContextI) (rps *readPools, err error) {

	var poolb []byte
	if poolb, err = ssc.getReadPoolsBytes(clientID, balances); err != nil {
		return
	}
	rps = newReadPools()
	err = rps.Decode(poolb)
	return
}

// newReadPool SC function creates new read pool for a client.
func (ssc *StorageSmartContract) newReadPool(t *transaction.Transaction,
	input []byte, balances chainState.StateContextI) (resp string, err error) {

	_, err = ssc.getReadPoolsBytes(t.ClientID, balances)

	if err != nil && err != util.ErrValueNotPresent {
		return "", common.NewError("new_read_pool_failed", err.Error())
	}

	if err == nil {
		return "", common.NewError("new_read_pool_failed", "already exist")
	}

	var rps = newReadPools()
	if err = rps.save(ssc.ID, t.ClientID, balances); err != nil {
		return "", common.NewError("new_read_pool_failed", err.Error())
	}

	return string(rps.Encode()), nil
}

func (ssc *StorageSmartContract) checkFill(t *transaction.Transaction,
	balances chainState.StateContextI) (err error) {

	var balance state.Balance
	balance, err = balances.GetClientBalance(t.ClientID)

	if err != nil && err != util.ErrValueNotPresent {
		return
	}

	if err == util.ErrValueNotPresent {
		return errors.New("no tokens to lock")
	}

	if state.Balance(t.Value) > balance {
		return errors.New("lock amount is greater than balance")
	}

	return
}

// lock tokens for read pool of transaction's client
func (ssc *StorageSmartContract) readPoolLock(t *transaction.Transaction,
	input []byte, balances chainState.StateContextI) (resp string, err error) {

	// configs

	var conf *readPoolConfig
	if conf, err = ssc.getReadPoolConfig(balances, true); err != nil {
		return "", common.NewError("read_pool_lock_failed",
			"can't get configs: "+err.Error())
	}

	// user read pools

	var rps *readPools
	if rps, err = ssc.getReadPools(t.ClientID, balances); err != nil {
		return "", common.NewError("read_pool_lock_failed", err.Error())
	}

	// lock request & user balance

	var lr lockRequest
	if err = lr.decode(input); err != nil {
		return "", common.NewError("read_pool_lock_failed", err.Error())
	}

	// check client balance
	if err = ssc.checkFill(t, balances); err != nil {
		return "", common.NewError("read_pool_lock_failed", err.Error())
	}

	// filter by configs

	if t.Value < conf.MinLock {
		return "", common.NewError("read_pool_lock_failed",
			"insufficient amount to lock")
	}

	if lr.Duration < conf.MinLockPeriod {
		return "", common.NewError("read_pool_lock_failed",
			fmt.Sprintf("duration (%s) is shorter than min lock period (%s)",
				lr.Duration.String(), conf.MinLockPeriod.String()))
	}

	if lr.Duration > conf.MaxLockPeriod {
		return "", common.NewError("read_pool_lock_failed",
			fmt.Sprintf("duration (%s) is longer than max lock period (%v)",
				lr.Duration.String(), conf.MaxLockPeriod.String()))
	}

	// lock

	var rp = newReadPool()
	rp.TokenLockInterface = &tokenLock{
		StartTime: t.CreationDate,
		Duration:  lr.Duration,
		Owner:     t.ClientID,
	}

	var transfer *state.Transfer
	if transfer, resp, err = rp.DigPool(t.Hash, t); err != nil {
		return "", common.NewError("read_pool_lock_failed",
			err.Error())
	}

	if err = balances.AddTransfer(transfer); err != nil {
		return "", common.NewError("read_pool_lock_failed", err.Error())
	}

	if err = rps.addPool(rp); err != nil {
		return "", common.NewError("read_pool_lock_failed", err.Error())
	}

	if err = rps.save(ssc.ID, t.ClientID, balances); err != nil {
		return "", common.NewError("read_pool_lock_failed", err.Error())
	}

	return
}

// unlock tokens if expired
func (ssc *StorageSmartContract) readPoolUnlock(t *transaction.Transaction,
	input []byte, balances chainState.StateContextI) (resp string, err error) {

	// user read pools

	var rps *readPools
	if rps, err = ssc.getReadPools(t.ClientID, balances); err != nil {
		return "", common.NewError("read_pool_unlock_failed", err.Error())
	}

	// the request

	var (
		transfer *state.Transfer
		req      unlockRequest
	)

	if err = req.decode(input); err != nil {
		return "", common.NewError("read_pool_unlock_failed", err.Error())
	}

	var pool, ok = rps.Pools[req.PoolID]
	if !ok {
		return "", common.NewError("read_pool_unlock_failed", "pool not found")
	}

	transfer, resp, err = pool.EmptyPool(ssc.ID, t.ClientID,
		common.ToTime(t.CreationDate))
	if err != nil {
		return "", common.NewError("read_pool_unlock_failed", err.Error())
	}

	if err = balances.AddTransfer(transfer); err != nil {
		return "", common.NewError("read_pool_unlock_failed", err.Error())
	}

	// save pools
	rps.delPool(pool.ID)
	if err = rps.save(ssc.ID, t.ClientID, balances); err != nil {
		return "", common.NewError("read_pool_unlock_failed", err.Error())
	}

	return
}

//
// stat
//

// statistic for all locked tokens of the read pool
func (ssc *StorageSmartContract) getReadPoolsStatsHandler(ctx context.Context,
	params url.Values, balances chainState.StateContextI) (
	resp interface{}, err error) {

	var (
		clientID = datastore.Key(params.Get("client_id"))
		rps      *readPools
	)
	if rps, err = ssc.getReadPools(clientID, balances); err != nil {
		return
	}

	if len(rps.Pools) == 0 {
		return nil, common.NewError("read_pool_stats", "no pools exist")
	}

	var (
		tp    = time.Now()
		stats readPoolStats
	)

	for _, rp := range rps.Pools {
		stat, err := rp.stat(tp)
		if err != nil {
			return nil, common.NewError("read_pool_stats", err.Error())
		}
		stats.addStat(stat)
	}

	return &stats, nil
}
