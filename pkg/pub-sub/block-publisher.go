package pubsub

import (
	"errors"
	"sync"

	evm "ethereum-parser/pkg/ethereum-rpc-client"

	"github.com/gammazero/deque"
)

/*
dev: This demo publisher including block storage and subscriber management.
*/

type BlockPublisher struct {
	sync.Mutex
	subs   map[*BlockSubscriber]bool
	blocks *deque.Deque[evm.Block]
}

var DefaultPublisher *BlockPublisher = NewBlockPublisher()

func SetDefaultPublisher(p *BlockPublisher) {
	DefaultPublisher = p
}

func NewBlockPublisher() *BlockPublisher {
	return &BlockPublisher{
		subs:   make(map[*BlockSubscriber]bool),
		blocks: deque.New[evm.Block](0),
	}
}

func (p *BlockPublisher) Subscribe(s *BlockSubscriber) error {
	p.Lock()
	p.subs[s] = true
	p.Unlock()

	return nil
}

func (p *BlockPublisher) Unsubscribe(s *BlockSubscriber) error {
	p.Lock()
	delete(p.subs, s)
	p.Unlock()

	return nil
}

func (p *BlockPublisher) AddBlock(block *evm.Block) error {
	p.Lock()
	p.blocks.PushBack(*block)
	p.Unlock()

	return nil
}

func (p *BlockPublisher) Publish(block *evm.Block) error {
	p.Lock()
	for s := range p.subs {
		s.Publish(block)
	}
	p.Unlock()

	return nil
}

func (p *BlockPublisher) GetLatestBlock() (*evm.Block, error) {
	p.Lock()
	defer p.Unlock()

	if p.blocks.Len() == 0 {
		return nil, errors.New("No blocks available")
	}

	blocks := p.blocks.Back()
	return &blocks, nil
}
