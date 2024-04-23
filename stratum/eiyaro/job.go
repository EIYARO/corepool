package ey

import (
	"math/big"
	"time"

	"corepool/core/protocol/bc"

	"corepool/common/mining/utils"
	ss "corepool/stratum"
	"corepool/stratum/eiyaro/util"
)

type eycJob struct {
	id                     ss.JobId
	version                uint64
	height                 uint64
	previousBlockHash      *bc.Hash
	timestamp              time.Time
	transactionsMerkleRoot *bc.Hash
	transactionStatusHash  *bc.Hash
	bits                   uint64
	seed                   *bc.Hash
	nonce                  uint64
	diff                   *big.Int
}

func (j *eycJob) GetId() ss.JobId {
	return j.id
}

func (j *eycJob) GetDiff() uint64 {
	return j.diff.Uint64()
}

func (j *eycJob) GetTarget() (string, bool, bool) {
	return "", false, false
}

func (j *eycJob) Encode() (interface{}, error) {
	return ss.StratumJSONRpcNotify{
		Version: "2.0",
		Method:  "job",
		Params:  j.genReplyData(),
	}, nil
}

func (j *eycJob) encodeLogin(login string) *jobReply {
	return &jobReply{
		Id:     login,
		Job:    j.genReplyData(),
		Status: "OK",
	}
}

func (j *eycJob) genReplyData() *jobReplyData {
	return &jobReplyData{
		JobId:                  j.GetId().String(),
		Version:                utils.ToLittleEndianHex(j.version),
		Height:                 utils.ToLittleEndianHex(j.height),
		PreviousBlockHash:      j.previousBlockHash.String(),
		Timestamp:              utils.ToLittleEndianHex(uint64(j.timestamp.Unix())),
		TransactionsMerkleRoot: j.transactionsMerkleRoot.String(),
		TransactionStatusHash:  j.transactionStatusHash.String(),
		Nonce:                  utils.ToLittleEndianHex(uint64(j.nonce)),
		Bits:                   utils.ToLittleEndianHex(j.bits),
		Seed:                   j.seed.String(),
		Target:                 util.GetTargetHex(j.diff),
	}
}
