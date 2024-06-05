package evmparser

import (
	"errors"
	evm "ethereum-parser/pkg/ethereum-rpc-client"
	"regexp"
	"strings"
	"sync"
)

type BasicEthereumParser struct {
	LastUpdateBlock int
	CurrentBlock    int
	Subscriptions   map[string]bool
	Transactions    map[string][]evm.Transaction
	mutex           sync.Mutex
}

func NewBasicEthereumParser() *BasicEthereumParser {
	block, err := evm.GetBlockNumber()
	if err != nil {
		return nil
	}

	return &BasicEthereumParser{
		LastUpdateBlock: block,
		CurrentBlock:    block,
		Subscriptions:   make(map[string]bool),
	}
}

func (p *BasicEthereumParser) GetCurrentBlock() (int, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	block, err := evm.GetBlockNumber()
	if err != nil {
		return -1, errors.New("error getting block number, " + err.Error())
	}
	p.CurrentBlock = block

	return p.CurrentBlock, nil
}

func (p *BasicEthereumParser) Subscribe(address string) (bool, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !isValidAddress(address) {
		return false, errors.New("invalid address")
	}

	if _, exists := p.Subscriptions[address]; exists {
		return false, nil
	}
	p.Subscriptions[strings.ToLower(address)] = true

	return true, nil
}

func (p *BasicEthereumParser) GetTransactions() ([]evm.Transaction, error) {
	blockNumber, err := evm.GetBlockNumber()
	if err != nil {
		return nil, errors.New("error getting block number, " + err.Error())
	}

	block, err := evm.GetBlockByNumber(blockNumber)
	if err != nil {
		return nil, errors.New("error getting block by number, " + err.Error())
	}

	var transactions []evm.Transaction
	for _, tx := range block.Transactions {
		if p.Subscriptions[strings.ToLower(tx.From)] || p.Subscriptions[strings.ToLower(tx.To)] {
			transactions = append(transactions, tx)
		}
	}

	return transactions, nil
}

func isValidAddress(v string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(v)
}
