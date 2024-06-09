package evmparser

import (
	"errors"
	evm "ethereum-parser/pkg/ethereum-rpc-client"
	"regexp"
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

func (p *BasicEthereumParser) GetCurrentBlock() (int, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	block, err := evm.GetBlockNumber()
	if err != nil {
		return -1, errors.New("error getting block number, " + err.Error())
	}

	return block, nil
}

func (p *BasicEthereumParser) Subscribe(address string) (bool, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !isValidAddress(address) {
		return false, errors.New("invalid address")
	}

	p.Subscriptions[strings.ToLower(address)] = true

	return true, nil
}

func (p *BasicEthereumParser) UnSubscribe(address string) (bool, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !isValidAddress(address) {
		return false, errors.New("invalid address")
	}

	p.Subscriptions[strings.ToLower(address)] = false

	return true, nil
}

func (p *BasicEthereumParser) GetTransactions() ([]evm.Transaction, error) {
	var transactions []evm.Transaction

	return transactions, nil
}

func isValidAddress(v string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(v)
}
