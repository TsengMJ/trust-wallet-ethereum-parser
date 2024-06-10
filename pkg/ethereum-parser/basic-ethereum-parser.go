package ethereumparser

import (
	"errors"
	evm "ethereum-parser/pkg/ethereum-rpc-client"
	pubsub "ethereum-parser/pkg/pub-sub"
	"ethereum-parser/util"
	"strings"
	"sync"
)

type BasicEthereumParser struct {
	Subscriptions map[string]bool
	mutex         sync.Mutex
}

func NewBasicEthereumParser() *BasicEthereumParser {
	return &BasicEthereumParser{
		Subscriptions: make(map[string]bool),
	}
}

func (p *BasicEthereumParser) GetCurrentBlock() (*evm.Block, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	block, err := pubsub.DefaultPublisher.GetLatestBlock()
	if err != nil {
		return nil, errors.New("error getting block number, " + err.Error())
	}

	return block, nil
}

func (p *BasicEthereumParser) Subscribe(address string) (bool, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !util.IsValidAddress(address) {
		return false, errors.New("invalid address")
	}

	p.Subscriptions[strings.ToLower(address)] = true

	return true, nil
}

func (p *BasicEthereumParser) UnSubscribe(address string) (bool, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !util.IsValidAddress(address) {
		return false, errors.New("invalid address")
	}

	p.Subscriptions[strings.ToLower(address)] = false

	return true, nil
}

func (p *BasicEthereumParser) GetTransactions() ([]evm.Transaction, error) {
	var transactions []evm.Transaction

	return transactions, nil
}
