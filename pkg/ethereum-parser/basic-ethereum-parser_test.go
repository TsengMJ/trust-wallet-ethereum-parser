// package evmparser_test

// import "testing"

// func TestNewBasicEthereumParser(t *testing.T) {

// }

package evmparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert" // using testify for assertions

	"ethereum-parser/config"
	evmparser "ethereum-parser/pkg/ethereum-parser" // replace with your package path (assuming ethereum-parser is a separate package
)

// Mock evm.GetBlockNumber function for testing
func mockGetBlockNumber(block int, err error) func() (int, error) {
	return func() (int, error) {
		return block, err
	}
}

func TestNewBasicEthereumParser(t *testing.T) {
	// Test Case 0: Success
	t.Run("successful creation", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Ethereum: config.Ethereum{
				Url: "https://eth.llamarpc.com",
			},
		}
		parser := evmparser.NewBasicEthereumParser()

		assert.NotNil(t, parser)
		assert.Empty(t, parser.Subscriptions)
	})

	// Test Case 1: Error getting block number
	t.Run("error getting block number", func(t *testing.T) {
		config.Config = config.EnvConfig{}
		parser := evmparser.NewBasicEthereumParser()
		assert.Nil(t, parser) // parser should be nil on error
	})
}

func TestGetCurrentBlock(t *testing.T) {
	t.Run("error getting block number", func(t *testing.T) {
		config.Config = config.EnvConfig{
			Ethereum: config.Ethereum{
				Url: "",
			},
		}
		parser := evmparser.BasicEthereumParser{
			Subscriptions: make(map[string]bool),
		}

		currentBlock, err := parser.GetCurrentBlock()

		assert.NotNil(t, err)
		assert.Equal(t, -1, currentBlock)
	})
}

func TestSubscribe(t *testing.T) {
	// Test Case 0: Success
	t.Run("successful subscription", func(t *testing.T) {
		parser := evmparser.BasicEthereumParser{
			Subscriptions: make(map[string]bool),
		}

		subscribed, err := parser.Subscribe("0x1234")
		assert.Nil(t, err)
		assert.True(t, subscribed)
		assert.Equal(t, map[string]bool{"0x1234": true}, parser.Subscriptions)
	})

	// Test Case 1: Duplicate subscription
	t.Run("duplicate subscription", func(t *testing.T) {
		parser := evmparser.BasicEthereumParser{
			Subscriptions: make(map[string]bool),
		}

		parser.Subscribe("0x1234")
		subscribed, err := parser.Subscribe("0x1234")

		assert.Nil(t, err)
		assert.False(t, subscribed)
		assert.Equal(t, map[string]bool{"0x1234": true}, parser.Subscriptions) // no change
	})

	// Test Case 2: Subscribe with uppercase address
	t.Run("subscribe with uppercase address", func(t *testing.T) {
		parser := evmparser.BasicEthereumParser{
			Subscriptions: make(map[string]bool),
		}

		subscribed, err := parser.Subscribe("0XABCD")

		assert.Nil(t, err)
		assert.True(t, subscribed)
		assert.Equal(t, map[string]bool{"0xabcd": true}, parser.Subscriptions)
	})
}

// Need a mock environment/implementation for testing GetTransactions
