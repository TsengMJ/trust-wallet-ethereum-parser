package util_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"ethereum-parser/util"
)

func TestIsValidAddress(t *testing.T) {
	tests := []struct {
		address  string
		expected bool
	}{
		{"0x123", false}, // Address too short
		{"0x7a250d5630B4cF539739dF2C5dAcb4c659F2488Dasdd", false}, // Address too long
		{"0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D", true},      // Valid address
		{"0x123g567890", false},                                   // Invalid character 'g'
	}

	for _, test := range tests {
		result := util.IsValidAddress(test.address)
		assert.Equal(t, test.expected, result, "Unexpected result for address: %s", test.address)
	}
}

func TestHexToDecimal(t *testing.T) {
	tests := []struct {
		hex         string
		expected    int64
		expectedErr error
	}{
		{"0x0", 0, nil},   // Hex 0
		{"0x1", 1, nil},   // Hex 1
		{"0x10", 16, nil}, // Hex 16
		{"0x7fffffffffffffff", 9223372036854775807, nil},             // Max int64 value
		{"0xFFFFFFFFFFFFFFFF", -1, errors.New("value out of range")}, // Overflow (max int64 + 1)
		{"0x123", 291, nil}, // Hex 291
	}

	for _, test := range tests {
		result, err := util.HexToDecimal(test.hex)

		if test.expectedErr != nil {
			assert.Error(t, err, "Expected error for hex: %s", test.hex)
			assert.ErrorContains(t, err, test.expectedErr.Error(), "Unexpected error for hex: %s", test.hex)
			continue
		} else {
			assert.NoError(t, err, "Unexpected error for hex: %s", test.hex)
			assert.Equal(t, test.expected, result, "Unexpected result for hex: %s", test.hex)
		}
	}
}
