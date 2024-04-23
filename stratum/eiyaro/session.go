package ey

import (
	"math"

	"github.com/cornelk/hashmap"

	ss "corepool/stratum"
)

const serverIdOffset = uint(60)

type eycSessionData struct {
	nonce    uint64
	worker   *ss.Worker
	submitId *hashmap.HashMap
}

func (s *eycSessionData) GetWorker() *ss.Worker {
	return s.worker
}

func (s *eycSessionData) SetWorker(worker *ss.Worker) {
	s.worker = worker
}

func (s *eycSessionData) getNonce() uint64 {
	return s.nonce
}

type eycSessionDataBuilder struct {
	id              uint64
	maxSessions     int
	sessionIdOffset uint
}

func NewBtmcSessionDataBuilder(serverId uint64, maxSessions int) *eycSessionDataBuilder {
	sessionIdOffset := serverIdOffset - uint(math.Ceil(math.Log2(float64(maxSessions))))
	if sessionIdOffset == serverIdOffset {
		sessionIdOffset--
	}
	return &eycSessionDataBuilder{
		id:              serverId,
		maxSessions:     maxSessions,
		sessionIdOffset: sessionIdOffset,
	}
}

// Build builds a eySession
func (b *eycSessionDataBuilder) Build(sessionId uint) ss.SessionData {
	return &eycSessionData{
		nonce:    (b.id << serverIdOffset) | (uint64(sessionId) << b.sessionIdOffset),
		worker:   nil,
		submitId: &hashmap.HashMap{},
	}
}
