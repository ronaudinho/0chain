package minersc

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"0chain.net/chaincore/block"
	"0chain.net/chaincore/state"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
	"0chain.net/core/datastore"
	"0chain.net/core/encryption"
	"0chain.net/core/util"

	"github.com/stretchr/testify/require"
)

//
// MPT based on memory-, level- and persistent node DB
//

type mptStore struct {
	mpt  util.MerklePatriciaTrieI
	mndb *util.MemoryNodeDB
	lndb *util.LevelNodeDB
	pndb *util.PNodeDB
	dir  string
}

func newMptStore(tb testing.TB) (mpts *mptStore) {
	mpts = new(mptStore)

	var dir, err = ioutil.TempDir("", "miner-sc-mpt")
	require.NoError(tb, err)

	mpts.pndb, err = util.NewPNodeDB(filepath.Join(dir, "data"),
		filepath.Join(dir, "log"))
	require.NoError(tb, err)

	mpts.merge(tb)

	mpts.dir = dir
	return
}

func (mpts *mptStore) Close() (err error) {
	if mpts == nil {
		return
	}
	if mpts.pndb != nil {
		mpts.pndb.Flush()
	}
	if mpts.dir != "" {
		err = os.RemoveAll(mpts.dir)
	}
	return
}

func (mpts *mptStore) merge(tb testing.TB) {
	if mpts == nil {
		return
	}

	var root util.Key

	if mpts.mndb != nil {
		root = mpts.mpt.GetRoot()
		require.NoError(tb, util.MergeState(
			context.Background(), mpts.mndb, mpts.pndb,
		))
		// mpts.pndb.Flush()
	}

	// for a worst case, no cached data, and we have to get everything from
	// the persistent store, from rocksdb

	mpts.mndb = util.NewMemoryNodeDB()                           //
	mpts.lndb = util.NewLevelNodeDB(mpts.mndb, mpts.pndb, false) // transaction
	mpts.mpt = util.NewMerklePatriciaTrie(mpts.lndb, 1)          //
	mpts.mpt.SetRoot(root)
}

//
// helper for tests implements chainState.StateContextI
//

type testBalances struct {
	balances      map[datastore.Key]state.Balance
	txn           *transaction.Transaction
	transfers     []*state.Transfer
	tree          map[datastore.Key]util.Serializable
	block         *block.Block
	blockSharders []string
	lfmb          *block.Block

	mpts      *mptStore
	skipMerge bool
}

func newTestBalances(t testing.TB, mpt bool) (tb *testBalances) {
	tb = &testBalances{
		balances: make(map[datastore.Key]state.Balance),
		tree:     make(map[datastore.Key]util.Serializable),
	}

	if mpt {
		tb.mpts = newMptStore(t)
	}

	return
}

func (tb *testBalances) setTransaction(t testing.TB,
	txn *transaction.Transaction) {

	tb.txn = txn

	if tb.mpts != nil && !tb.skipMerge {
		tb.mpts.merge(t)
	}
}

func (tb *testBalances) setBalance(key datastore.Key, b state.Balance) {
	tb.balances[key] = b
}

func (tb *testBalances) setLFMB(lfmb *block.Block) {
	tb.lfmb = lfmb
}

func (tb *testBalances) GetBlock() *block.Block {
	return tb.block
}

func (tb *testBalances) SetMagicBlock(mb *block.MagicBlock) {
	if tb.block != nil {
		tb.block.MagicBlock = mb
	}
}

func (tb *testBalances) GetBlockSharders(*block.Block) []string {
	return tb.blockSharders
}

// stubs
func (tb *testBalances) GetState() util.MerklePatriciaTrieI       { return nil }
func (tb *testBalances) GetTransaction() *transaction.Transaction { return nil }
func (tb *testBalances) Validate() error                          { return nil }
func (tb *testBalances) GetMints() []*state.Mint                  { return nil }
func (tb *testBalances) SetStateContext(*state.State) error       { return nil }
func (tb *testBalances) GetTransfers() []*state.Transfer          { return nil }
func (tb *testBalances) AddSignedTransfer(st *state.SignedTransfer) {
}
func (tb *testBalances) GetSignedTransfers() []*state.SignedTransfer {
	return nil
}

func (tb *testBalances) DeleteTrieNode(key datastore.Key) (
	datastore.Key, error) {

	if tb.mpts != nil {
		if encryption.IsHash(key) {
			return "", common.NewError("failed to get trie node",
				"key is too short")
		}
		var btkey, err = tb.mpts.mpt.Delete(util.Path(encryption.Hash(key)))
		return datastore.Key(btkey), err
	}

	delete(tb.tree, key)
	return "", nil
}

func (tb *testBalances) GetLastestFinalizedMagicBlock() *block.Block {
	return tb.lfmb
}

func (tb *testBalances) GetSignatureScheme() encryption.SignatureScheme {
	return encryption.NewBLS0ChainScheme()
}

func (tb *testBalances) GetClientBalance(clientID datastore.Key) (
	b state.Balance, err error) {

	var ok bool
	if b, ok = tb.balances[clientID]; !ok {
		return 0, util.ErrValueNotPresent
	}
	return
}

func (tb *testBalances) GetTrieNode(key datastore.Key) (
	node util.Serializable, err error) {

	if encryption.IsHash(key) {
		return nil, common.NewError("failed to get trie node",
			"key is too short")
	}

	if tb.mpts != nil {
		return tb.mpts.mpt.GetNodeValue(util.Path(encryption.Hash(key)))
	}

	var ok bool
	if node, ok = tb.tree[key]; !ok {
		return nil, util.ErrValueNotPresent
	}
	return
}

func (tb *testBalances) InsertTrieNode(key datastore.Key,
	node util.Serializable) (datastore.Key, error) {

	if tb.mpts != nil {
		if encryption.IsHash(key) {
			return "", common.NewError("failed to get trie node",
				"key is too short")
		}
		var btkey, err = tb.mpts.mpt.Insert(util.Path(encryption.Hash(key)), node)
		return datastore.Key(btkey), err
	}

	tb.tree[key] = node
	return "", nil
}

func (tb *testBalances) AddTransfer(t *state.Transfer) error {
	if t.ClientID != tb.txn.ClientID && t.ClientID != tb.txn.ToClientID {
		return state.ErrInvalidTransfer
	}
	tb.balances[t.ClientID] -= t.Amount
	tb.balances[t.ToClientID] += t.Amount
	tb.transfers = append(tb.transfers, t)
	return nil
}

func (tb *testBalances) AddMint(mint *state.Mint) error {
	if mint.Minter != ADDRESS {
		panic("invalid miner: " + mint.Minter)
	}
	tb.balances[mint.ToClientID] += mint.Amount // mint!
	return nil
}
