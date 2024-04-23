package ey

import (
	"errors"
	"sync"

	"github.com/segmentio/encoding/json"

	"corepool/common/logger"
	"corepool/core/api"
	ss "corepool/stratum"
	"corepool/stratum/eiyaro/rpc"
)

type eycNodeSyncer struct {
	client *rpc.BtmcClient
	bt     *api.GetWorkResp
	btLock sync.RWMutex

	latestHeight uint64
}

func NewBtmcNodeSyncer(service string, nodeURL string) (*eycNodeSyncer, error) {
	return &eycNodeSyncer{
		client:       rpc.NewBtmcClient(service, nodeURL),
		latestHeight: 0,
	}, nil
}

func (n *eycNodeSyncer) fetchBlockTemplate() (ss.BlockTemplate, error) {
	reply, err := n.client.GetWork()
	if err != nil {
		return nil, err
	}

	header := reply.BlockHeader
	if header == nil {
		return nil, ErrNullBlockHeader
	}

	return &eiyaroBlockTemplate{
		version:                header.Version,
		height:                 header.Height,
		previousBlockHash:      &header.PreviousBlockHash,
		timestamp:              header.Time(),
		transactionsMerkleRoot: &header.TransactionsMerkleRoot,
		transactionStatusHash:  &header.TransactionStatusHash,
		nonce:                  header.Nonce,
		bits:                   header.Bits,
		seed:                   reply.Seed,
	}, nil
}

func (n *eycNodeSyncer) Pull() (ss.BlockTemplate, error) {
	return n.fetchBlockTemplate()
}

func (n *eycNodeSyncer) Submit(share ss.Share) error {
	eycShare := share.(*eycShare)
	rawdata, err := n.client.SubmitBlock(&api.SubmitWorkReq{BlockHeader: eycShare.header})
	if err != nil {
		return err
	}

	resultrawdata, err := json.Marshal(rawdata)
	if err != nil {
		return err
	}
	var result bool
	if err := json.Unmarshal(resultrawdata, &result); err != nil {
		return err
	}
	if !result {
		logger.Error("block rejected", "nonce", eycShare.nonce, "hash", eycShare.blockHash)
		return nil
	}
	logger.Info("send nonce success", "nonce", eycShare.nonce)
	return nil
}

func (n *eycNodeSyncer) GetBt() (*api.GetWorkResp, error) {
	n.btLock.RLock()
	defer n.btLock.RUnlock()
	if n.bt == nil {
		return nil, errors.New("getting blocktemplate")
	}
	return n.bt, nil
}
