package ethereumparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	evmparser "ethereum-parser/pkg/ethereum-parser"
)

func TestBasicEthereumParser(t *testing.T) {
	t.Run("GetCurrentBlock", func(t *testing.T) {
		parser := evmparser.NewBasicEthereumParser()

		block, err := parser.GetCurrentBlock()

		assert.Nil(t, block, "Block should be nil")
		assert.Error(t, err, "Error should occur")
		assert.ErrorContains(t, err, "error getting block number,")
	})

	t.Run("Subscribe", func(t *testing.T) {
		parser := evmparser.NewBasicEthereumParser()

		// Test with invalid address
		valid, err := parser.Subscribe("invalid_address")

		assert.False(t, valid, "Should return false for invalid address")
		assert.Error(t, err, "Error should occur")
		assert.EqualError(t, err, "invalid address")

		// Test with valid address
		valid, err = parser.Subscribe("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D")

		assert.True(t, valid, "Should return true for valid address")
		assert.NoError(t, err, "No error expected")
	})

	t.Run("UnSubscribe", func(t *testing.T) {
		parser := evmparser.NewBasicEthereumParser()

		// Test with invalid address
		valid, err := parser.UnSubscribe("invalid_address")

		assert.False(t, valid, "Should return false for invalid address")
		assert.Error(t, err, "Error should occur")
		assert.EqualError(t, err, "invalid address")

		// Test with valid address
		valid, err = parser.UnSubscribe("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D")

		assert.True(t, valid, "Should return true for valid address")
		assert.NoError(t, err, "No error expected")
	})

	t.Run("GetTransactions", func(t *testing.T) {
		parser := evmparser.NewBasicEthereumParser()

		transactions, err := parser.GetTransactions()

		assert.Empty(t, transactions, "Transactions should be empty")
		assert.NoError(t, err, "No error expected")
	})
}
