package ey

import (
	"math/big"

	"corepool/core/consensus/difficulty"
	algorithm "corepool/core/mining/tensority/go_algorithm"
	"corepool/core/protocol/bc/types"

	ss "corepool/stratum"
)

type eycVerifier struct {
	serverState *ss.ServerState
}

func NewBtmcVerifier(state *ss.ServerState) (*eycVerifier, error) {
	return &eycVerifier{
		serverState: state,
	}, nil
}

func (v *eycVerifier) Verify(share ss.Share) error {
	eycShare := share.(*eycShare)
	eycJob := eycShare.job
	eycShare.header = &types.BlockHeader{
		Version:           eycJob.version,
		Height:            eycJob.height,
		PreviousBlockHash: *eycJob.previousBlockHash,
		Timestamp:         uint64(eycJob.timestamp.Unix()),
		BlockCommitment: types.BlockCommitment{
			TransactionsMerkleRoot: *eycJob.transactionsMerkleRoot,
			TransactionStatusHash:  *eycJob.transactionStatusHash,
		},
		Nonce: eycShare.nonce,
		Bits:  eycJob.bits,
	}
	shareHeader := eycShare.header
	headerHash := shareHeader.Hash()
	cmpHash := algorithm.LegacyAlgorithm(&headerHash, eycJob.seed)
	if cmpHash == nil {
		share.UpdateState(ss.ShareStateRejected, ss.RejectReasonUndefined)
		return nil
	}

	eycShare.blockHash = &headerHash
	bMax := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil)
	bits := difficulty.BigToCompact(big.NewInt(0).Div(bMax, eycShare.netDiff))
	if difficulty.HashToBig(cmpHash).Cmp(difficulty.CompactToBig(bits)) <= 0 {
		share.UpdateState(ss.ShareStateBlock, ss.RejectReasonPass)
		return nil
	}

	shareBits := difficulty.BigToCompact(big.NewInt(0).Div(bMax, eycJob.diff))
	if difficulty.HashToBig(cmpHash).Cmp(difficulty.CompactToBig(shareBits)) > 0 {
		share.UpdateState(ss.ShareStateRejected, ss.RejectReasonLowDiff)
		return nil
	}
	share.UpdateState(ss.ShareStateAccepted, ss.RejectReasonPass)
	return nil
}
