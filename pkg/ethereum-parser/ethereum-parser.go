package evmparser

import (
	evm "ethereum-parser/pkg/ethereum-rpc-client"
)

type EthereumParser interface {
	GetCurrentBlock() (int, error)
	Subscribe(address string) (bool, error)
	UnSubscribe(address string) (bool, error)
	GetTransactions() ([]evm.Transaction, error)
}
