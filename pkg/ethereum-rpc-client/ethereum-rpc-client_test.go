package ethereumrpcclient_test

import (
	"ethereum-parser/config"
	evmrpcclient "ethereum-parser/pkg/ethereum-rpc-client"
	"testing"
)

func TestCallJSONRPC(t *testing.T) {
	// Test Case 0: Success
	t.Run("Success", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Ethereum: config.Ethereum{
				Url: "https://eth.llamarpc.com",
			},
		}

		result, err := evmrpcclient.CallJSONRPC("eth_blockNumber", []interface{}{})
		if err != nil {
			t.Fatal("Expected no error, got ", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}
	})

	// Test Case 1: Invalid method
	t.Run("Invalid method", func(t *testing.T) {
		result, err := evmrpcclient.CallJSONRPC("", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if result != nil {
			t.Fatal("Expected nil, got ", result)
		}
	})

	// Test Case 2: Invalid url
	t.Run("Invalid url", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Ethereum: config.Ethereum{
				Url: "",
			},
		}

		result, err := evmrpcclient.CallJSONRPC("eth_blockNumber", []interface{}{})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if result != nil {
			t.Fatal("Expected nil, got ", result)
		}
	})

	// Test Case 3: Invalid request body
	t.Run("Invalid request body", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Ethereum: config.Ethereum{
				Url: "https://eth.llamarpc.com",
			},
		}

		result, err := evmrpcclient.CallJSONRPC("eth_getTransactionCount", nil)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if result != nil {
			t.Fatal("Expected nil, got ", result)
		}
	})
}

func TestGetBlockNumber(t *testing.T) {
	// Test Case 0: Success
	t.Run("Success", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Ethereum: config.Ethereum{
				Url: "https://eth.llamarpc.com",
			},
		}

		blockNumber, err := evmrpcclient.GetBlockNumber()
		if err != nil {
			t.Fatal("Expected no error, got ", err)
		}

		if blockNumber == 0 {
			t.Fatal("Expected block number, got 0")
		}
	})

	// Test Case 1: Invalid url
	t.Run("Invalid url", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Ethereum: config.Ethereum{
				Url: "",
			},
		}

		blockNumber, err := evmrpcclient.GetBlockNumber()
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if blockNumber != 0 {
			t.Fatal("Expected 0, got ", blockNumber)
		}
	})
}

func TestGetBlockByNumber(t *testing.T) {
	// Test Case 0: Success
	t.Run("Success", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Ethereum: config.Ethereum{
				Url: "https://eth.llamarpc.com",
			},
		}

		block, err := evmrpcclient.GetBlockByNumber(1)
		if err != nil {
			t.Fatal("Expected no error, got ", err)
		}

		if block == nil {
			t.Fatal("Expected block, got nil")
		}
	})

	// Test Case 1: Invalid url
	t.Run("Invalid url", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Ethereum: config.Ethereum{
				Url: "",
			},
		}

		block, err := evmrpcclient.GetBlockByNumber(1)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if block != nil {
			t.Fatal("Expected nil, got ", block)
		}
	})

	// Test Case 2: Invalid block number
	t.Run("Invalid block number", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Ethereum: config.Ethereum{
				Url: "https://eth.llamarpc.com",
			},
		}

		block, err := evmrpcclient.GetBlockByNumber(-1)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if block != nil {
			t.Fatal("Expected nil, got ", block)
		}
	})

}
