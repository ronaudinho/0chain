// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	block "0chain.net/chaincore/block"
	chain "0chain.net/chaincore/chain"

	context "context"

	datastore "0chain.net/core/datastore"

	miner "0chain.net/miner"

	mock "github.com/stretchr/testify/mock"

	round "0chain.net/chaincore/round"
)

// Protocol is an autogenerated mock type for the Protocol type
type Protocol struct {
	mock.Mock
}

// AddToRoundVerification provides a mock function with given fields: ctx, r, b
func (_m *Protocol) AddToRoundVerification(ctx context.Context, r *miner.Round, b *block.Block) {
	_m.Called(ctx, r, b)
}

// AddVRFShare provides a mock function with given fields: ctx, r, vrfs
func (_m *Protocol) AddVRFShare(ctx context.Context, r *miner.Round, vrfs *round.VRFShare) bool {
	ret := _m.Called(ctx, r, vrfs)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, *miner.Round, *round.VRFShare) bool); ok {
		r0 = rf(ctx, r, vrfs)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// AddVerificationTicket provides a mock function with given fields: ctx, b, bvt
func (_m *Protocol) AddVerificationTicket(ctx context.Context, b *block.Block, bvt *block.VerificationTicket) bool {
	ret := _m.Called(ctx, b, bvt)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, *block.Block, *block.VerificationTicket) bool); ok {
		r0 = rf(ctx, b, bvt)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// CancelRoundVerification provides a mock function with given fields: ctx, r
func (_m *Protocol) CancelRoundVerification(ctx context.Context, r *miner.Round) {
	_m.Called(ctx, r)
}

// CollectBlocksForVerification provides a mock function with given fields: ctx, r
func (_m *Protocol) CollectBlocksForVerification(ctx context.Context, r *miner.Round) {
	_m.Called(ctx, r)
}

// FinalizeBlock provides a mock function with given fields: ctx, b
func (_m *Protocol) FinalizeBlock(ctx context.Context, b *block.Block) error {
	ret := _m.Called(ctx, b)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *block.Block) error); ok {
		r0 = rf(ctx, b)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FinalizeRound provides a mock function with given fields: ctx, r, bsh
func (_m *Protocol) FinalizeRound(ctx context.Context, r round.RoundI, bsh chain.BlockStateHandler) {
	_m.Called(ctx, r, bsh)
}

// GenerateBlock provides a mock function with given fields: ctx, b, bsh, waitOver
func (_m *Protocol) GenerateBlock(ctx context.Context, b *block.Block, bsh chain.BlockStateHandler, waitOver bool) error {
	ret := _m.Called(ctx, b, bsh, waitOver)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *block.Block, chain.BlockStateHandler, bool) error); ok {
		r0 = rf(ctx, b, bsh, waitOver)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HandleNotarizationMessage provides a mock function with given fields: ctx, msg
func (_m *Protocol) HandleNotarizationMessage(ctx context.Context, msg *miner.BlockMessage) {
	_m.Called(ctx, msg)
}

// HandleNotarizedBlockMessage provides a mock function with given fields: ctx, msg
func (_m *Protocol) HandleNotarizedBlockMessage(ctx context.Context, msg *miner.BlockMessage) {
	_m.Called(ctx, msg)
}

// HandleRoundTimeout provides a mock function with given fields: ctx
func (_m *Protocol) HandleRoundTimeout(ctx context.Context) {
	_m.Called(ctx)
}

// HandleVRFShare provides a mock function with given fields: ctx, msg
func (_m *Protocol) HandleVRFShare(ctx context.Context, msg *miner.BlockMessage) {
	_m.Called(ctx, msg)
}

// HandleVerificationTicketMessage provides a mock function with given fields: ctx, msg
func (_m *Protocol) HandleVerificationTicketMessage(ctx context.Context, msg *miner.BlockMessage) {
	_m.Called(ctx, msg)
}

// HandleVerifyBlockMessage provides a mock function with given fields: ctx, msg
func (_m *Protocol) HandleVerifyBlockMessage(ctx context.Context, msg *miner.BlockMessage) {
	_m.Called(ctx, msg)
}

// IsBlockNotarized provides a mock function with given fields: ctx, b
func (_m *Protocol) IsBlockNotarized(ctx context.Context, b *block.Block) bool {
	ret := _m.Called(ctx, b)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, *block.Block) bool); ok {
		r0 = rf(ctx, b)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ProcessVerifiedTicket provides a mock function with given fields: ctx, r, b, vt
func (_m *Protocol) ProcessVerifiedTicket(ctx context.Context, r *miner.Round, b *block.Block, vt *block.VerificationTicket) {
	_m.Called(ctx, r, b, vt)
}

// SaveMagicBlock provides a mock function with given fields:
func (_m *Protocol) SaveMagicBlock() chain.MagicBlockSaveFunc {
	ret := _m.Called()

	var r0 chain.MagicBlockSaveFunc
	if rf, ok := ret.Get(0).(func() chain.MagicBlockSaveFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chain.MagicBlockSaveFunc)
		}
	}

	return r0
}

// SendBlock provides a mock function with given fields: ctx, b
func (_m *Protocol) SendBlock(ctx context.Context, b *block.Block) {
	_m.Called(ctx, b)
}

// SendFinalizedBlock provides a mock function with given fields: ctx, b
func (_m *Protocol) SendFinalizedBlock(ctx context.Context, b *block.Block) {
	_m.Called(ctx, b)
}

// SendNotarization provides a mock function with given fields: ctx, b
func (_m *Protocol) SendNotarization(ctx context.Context, b *block.Block) {
	_m.Called(ctx, b)
}

// SendNotarizedBlock provides a mock function with given fields: ctx, b
func (_m *Protocol) SendNotarizedBlock(ctx context.Context, b *block.Block) {
	_m.Called(ctx, b)
}

// SendVRFShare provides a mock function with given fields: ctx, r
func (_m *Protocol) SendVRFShare(ctx context.Context, r *round.VRFShare) {
	_m.Called(ctx, r)
}

// SendVerificationTicket provides a mock function with given fields: ctx, b, bvt
func (_m *Protocol) SendVerificationTicket(ctx context.Context, b *block.Block, bvt *block.BlockVerificationTicket) {
	_m.Called(ctx, b, bvt)
}

// StartNextRound provides a mock function with given fields: ctx, _a1
func (_m *Protocol) StartNextRound(ctx context.Context, _a1 *miner.Round) *miner.Round {
	ret := _m.Called(ctx, _a1)

	var r0 *miner.Round
	if rf, ok := ret.Get(0).(func(context.Context, *miner.Round) *miner.Round); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*miner.Round)
		}
	}

	return r0
}

// UpdateFinalizedBlock provides a mock function with given fields: ctx, b
func (_m *Protocol) UpdateFinalizedBlock(ctx context.Context, b *block.Block) {
	_m.Called(ctx, b)
}

// UpdatePendingBlock provides a mock function with given fields: ctx, b, txns
func (_m *Protocol) UpdatePendingBlock(ctx context.Context, b *block.Block, txns []datastore.Entity) {
	_m.Called(ctx, b, txns)
}

// ValidateMagicBlock provides a mock function with given fields: _a0, _a1, _a2
func (_m *Protocol) ValidateMagicBlock(_a0 context.Context, _a1 *round.Round, _a2 *block.Block) bool {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, *round.Round, *block.Block) bool); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// VerifyBlock provides a mock function with given fields: ctx, b
func (_m *Protocol) VerifyBlock(ctx context.Context, b *block.Block) (*block.BlockVerificationTicket, error) {
	ret := _m.Called(ctx, b)

	var r0 *block.BlockVerificationTicket
	if rf, ok := ret.Get(0).(func(context.Context, *block.Block) *block.BlockVerificationTicket); ok {
		r0 = rf(ctx, b)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*block.BlockVerificationTicket)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *block.Block) error); ok {
		r1 = rf(ctx, b)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyNotarization provides a mock function with given fields: ctx, b, bvt, _a3
func (_m *Protocol) VerifyNotarization(ctx context.Context, b *block.Block, bvt []*block.VerificationTicket, _a3 int64) error {
	ret := _m.Called(ctx, b, bvt, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *block.Block, []*block.VerificationTicket, int64) error); ok {
		r0 = rf(ctx, b, bvt, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyTicket provides a mock function with given fields: ctx, blockHash, vt, _a3
func (_m *Protocol) VerifyTicket(ctx context.Context, blockHash string, vt *block.VerificationTicket, _a3 int64) error {
	ret := _m.Called(ctx, blockHash, vt, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *block.VerificationTicket, int64) error); ok {
		r0 = rf(ctx, blockHash, vt, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}