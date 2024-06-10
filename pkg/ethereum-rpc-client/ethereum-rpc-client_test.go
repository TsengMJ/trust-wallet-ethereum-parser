package ethereumrpcclient_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"ethereum-parser/config"
	ethereumrpcclient "ethereum-parser/pkg/ethereum-rpc-client"
)

func TestCallJSONRPC(t *testing.T) {
	cases := []struct {
		name   string
		method string
		params []interface{}

		responseStatus int
		responseBody   string
		expected       json.RawMessage
		expectedErr    error
		expectedPanic  bool
	}{
		{
			name:           "Valid response",
			method:         "eth_blockNumber",
			params:         []interface{}{},
			responseStatus: http.StatusOK,
			responseBody:   `{"jsonrpc":"2.0","result":"0x1","id":1}`,
			expected:       json.RawMessage(`"0x1"`),
			expectedErr:    nil,
		},
		{
			name:           "Empty method",
			method:         "",
			params:         []interface{}{},
			responseStatus: http.StatusOK,
			responseBody:   "",
			expected:       nil,
			expectedErr:    errors.New("method is empty"),
			expectedPanic:  false,
		},
		{
			name:           "Empty Ethereum URL",
			method:         "eth_blockNumber",
			params:         []interface{}{},
			responseStatus: 0,
			responseBody:   "",
			expected:       nil,
			expectedErr:    errors.New("ethereum url is empty"),
		},
		{
			name:           "Error sending HTTP request",
			method:         "eth_blockNumber",
			params:         []interface{}{},
			responseStatus: http.StatusBadGateway,
			responseBody:   "",
			expected:       nil,
			expectedErr:    errors.New("error sending request"),
		},
		{
			name:           "Error decoding response body",
			method:         "eth_blockNumber",
			params:         []interface{}{},
			responseStatus: http.StatusOK,
			responseBody:   "invalid json",
			expected:       nil,
			expectedErr:    errors.New("error decoding response body"),
		},
		{
			name:           "Error getting block number",
			method:         "eth_blockNumber",
			params:         []interface{}{},
			responseStatus: http.StatusOK,
			responseBody:   `{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":1}`,
			expected:       nil,
			expectedErr:    errors.New("Method not found"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Set up a mock HTTP server
			config.Config.Ethereum.Url = ""

			if c.responseStatus != 0 {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(c.responseStatus)
					w.Write([]byte(c.responseBody))
				}))
				defer server.Close()

				// Set Ethereum URL in configuration
				config.Config.Ethereum.Url = server.URL
			}

			// Call CallJSONRPC and capture panic if any
			defer func() {
				if r := recover(); r != nil {
					assert.True(t, c.expectedPanic, "Unexpected panic occurred")
				}
			}()

			result, err := ethereumrpcclient.CallJSONRPC(c.method, c.params)

			if c.expectedErr != nil {
				assert.Error(t, err, "Expected an error")
				assert.ErrorContains(t, err, c.expectedErr.Error(), "Unexpected error")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, c.expected, result, "Result does not match expected")
			}
		})
	}
}

func TestGetBlockNumber(t *testing.T) {
	cases := []struct {
		name          string
		response      string
		expected      int
		expectedErr   error
		expectedPanic bool
	}{
		{
			name:        "Valid block number",
			response:    `{"jsonrpc":"2.0","result":"0x1","id":1}`,
			expected:    1,
			expectedErr: nil,
		},
		{
			name:          "Error getting block number",
			response:      `{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":1}`,
			expected:      0,
			expectedErr:   errors.New("error getting block number, Method not found"),
			expectedPanic: false,
		},
		{
			name:          "Empty result",
			response:      `{"jsonrpc":"2.0","result":"","id":1}`,
			expected:      0,
			expectedErr:   errors.New("error unmarshalling block number"),
			expectedPanic: true,
		},
		{
			name:          "Invalid response format",
			response:      `"invalid json"`,
			expected:      0,
			expectedErr:   errors.New("error decoding response body"),
			expectedPanic: false,
		},
		{
			name:          "Error decoding response body",
			response:      `"invalid json"`,
			expected:      0,
			expectedErr:   errors.New("error decoding response body"),
			expectedPanic: false,
		},
		{
			name:          "Error getting block number",
			response:      `{"jsonrpc":"2.0","error":{"code":-32000,"message":"Internal error"},"id":1}`,
			expected:      0,
			expectedErr:   errors.New("error getting block number, Internal error"),
			expectedPanic: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Set up a mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(c.response))
			}))
			defer server.Close()

			// Set Ethereum URL in configuration
			config.Config.Ethereum.Url = server.URL

			// Call GetBlockNumber and capture panic if any
			defer func() {
				if r := recover(); r != nil {
					assert.True(t, c.expectedPanic, "Unexpected panic occurred")
				}
			}()

			blockNumber, err := ethereumrpcclient.GetBlockNumber()

			if c.expectedErr != nil {
				assert.Error(t, err, "Expected an error")
				assert.ErrorContains(t, err, c.expectedErr.Error(), "Unexpected error")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, c.expected, blockNumber, "Block number does not match expected")
			}
		})
	}
}

func TestGetBlockByNumber(t *testing.T) {
	cases := []struct {
		name          string
		blockNumber   int
		response      string
		expected      *ethereumrpcclient.Block
		expectedErr   error
		expectedPanic bool
	}{
		{
			name:        "Valid block",
			blockNumber: 1,
			response:    `{"jsonrpc":"2.0","result":{"number":"0x1","hash":"0x123","parentHash":"0x456","miner":"0x789"},"id":1}`,
			expected: &ethereumrpcclient.Block{
				Number:     "0x1",
				Hash:       "0x123",
				ParentHash: "0x456",
				Miner:      "0x789",
			},
			expectedErr:   nil,
			expectedPanic: false,
		},
		{
			name:          "Error getting block",
			blockNumber:   1,
			response:      `{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":1}`,
			expected:      nil,
			expectedErr:   errors.New("error getting block, Method not found"),
			expectedPanic: false,
		},
		{
			name:          "Empty result",
			blockNumber:   1,
			response:      `{"jsonrpc":"2.0","result":"","id":1}`,
			expected:      nil,
			expectedErr:   errors.New("error unmarshalling block"),
			expectedPanic: false,
		},
		{
			name:          "Invalid response format",
			blockNumber:   1,
			response:      `"invalid json"`,
			expected:      nil,
			expectedErr:   errors.New("error decoding response body"),
			expectedPanic: false,
		},
		{
			name:        "Valid block with transactions",
			blockNumber: 1,
			response:    `{"jsonrpc":"2.0","result":{"number":"0x1","hash":"0x123","transactions":[{"hash":"0xabc","from":"0xdef","to":"0x456","value":"100"}]},"id":1}`,
			expected: &ethereumrpcclient.Block{
				Number: "0x1",
				Hash:   "0x123",
				Transactions: []ethereumrpcclient.Transaction{
					{
						Hash:  "0xabc",
						From:  "0xdef",
						To:    "0x456",
						Value: "100",
					},
				},
			},
			expectedErr:   nil,
			expectedPanic: false,
		},
		{
			name:          "Empty result transactions",
			blockNumber:   1,
			response:      `{"jsonrpc":"2.0","result":{"number":"0x1","hash":"0x123","transactions":[]},"id":1}`,
			expected:      &ethereumrpcclient.Block{Number: "0x1", Hash: "0x123", Transactions: []ethereumrpcclient.Transaction{}},
			expectedErr:   nil,
			expectedPanic: false,
		},
		{
			name:          "Error getting block",
			blockNumber:   1,
			response:      `{"jsonrpc":"2.0","error":{"code":-32000,"message":"Internal error"},"id":1}`,
			expected:      nil,
			expectedErr:   errors.New("error getting block, Internal error"),
			expectedPanic: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Set up a mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(c.response))
			}))
			defer server.Close()

			// Set Ethereum URL in configuration
			config.Config.Ethereum.Url = server.URL

			// Call GetBlockByNumber and capture panic if any
			defer func() {
				if r := recover(); r != nil {
					assert.True(t, c.expectedPanic, "Unexpected panic occurred")
				}
			}()

			block, err := ethereumrpcclient.GetBlockByNumber(c.blockNumber)

			if c.expectedErr != nil {
				assert.Error(t, err, "Expected an error")
				assert.ErrorContains(t, err, c.expectedErr.Error(), "Unexpected error")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, c.expected, block, "Block does not match expected")
			}
		})
	}
}
