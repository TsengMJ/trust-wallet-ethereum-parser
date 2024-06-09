package pubsub

import (
	evm "ethereum-parser/pkg/ethereum-rpc-client"
	"sync"
)

type BlockSubscriber struct {
	sync.Mutex
	Handler chan *evm.Block
	Quit    chan struct{}
}

func NewBlockSubscriber() *BlockSubscriber {
	return &BlockSubscriber{
		Handler: make(chan *evm.Block, 1024),
		Quit:    make(chan struct{}),
	}
}

func (s *BlockSubscriber) Publish(block *evm.Block) {
	select {
	case s.Handler <- block:
	default:
	}
}
