package rpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/segmentio/encoding/json"

	"corepool/core/api"
	eyrpc "corepool/core/blockchain/rpc"
	ss "corepool/stratum"
)

type EyClient struct {
	service string
	eyCli *eyrpc.Client
}

func NewEyClient(service string, nodeURL string) *EyClient {
	return &EyClient{
		service: service,
		eyCli: &eyrpc.Client{BaseURL: nodeURL},
	}
}

func (c *EyClient) GetWork() (*api.GetWorkResp, error) {
	var result ss.NodeJsonRpcResp
	if err := ss.CallWithMethod(c.service, "get-work", []string{}, &result); err != nil {
		return nil, err
	}

	if result.Data == nil {
		return nil, errors.New("empty reply for get work")
	}
	var reply api.GetWorkResp
	if err := json.Unmarshal(*result.Data, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (c *EyClient) SubmitBlock(req interface{}) (interface{}, error) {
	var response = &api.Response{}
	c.eyCli.Call(context.Background(), "/submit-work", req, response)

	switch response.Status {
	case api.FAIL:
		return nil, errors.New("error.eiyarod")
	case "":
		return nil, errors.New("error.connect_eiyarod")
	}
	return response.Data, nil
}

type Peer struct {
	RemoteAddr string `json:"remote_addr"`
	Height     int64  `json:"height"`
	Ping       string `json:"ping"`
}

func (c *EyClient) GetPeers() ([]*Peer, error) {
	var result ss.NodeJsonRpcResp
	if err := ss.CallWithMethod(c.service, "list-peers", []string{}, &result); err != nil {
		return nil, err
	}

	var peers []*Peer
	if err := json.Unmarshal(*result.Data, &peers); err != nil {
		return nil, err
	}

	return peers, nil
}

type balance struct {
	Amount int64 `json:"amount"`
}

func (c *EyClient) GetBalance() (int64, error) {
	var result ss.NodeJsonRpcResp
	if err := ss.CallWithMethod(c.service, "list-balances", []string{}, &result); err != nil {
		return 0, err
	}

	var balances []*balance
	if err := json.Unmarshal(*result.Data, &balances); err != nil {
		return 0, err
	}

	if len(balances) != 1 {
		return 0, fmt.Errorf("unexpected balance response")
	}

	return balances[0].Amount, nil
}
