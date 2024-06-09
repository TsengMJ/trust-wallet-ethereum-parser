package ethereumrpcclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"ethereum-parser/config"
	"net/http"
	"strconv"
)

type Transaction struct {
	Hash  string
	From  string
	To    string
	Value string
}

type Block struct {
	Number           string
	Hash             string
	ParentHash       string
	Nonce            string
	Sha3Uncles       string
	LogsBloom        string
	TransactionsRoot string
	StateRoot        string
	ReceiptsRoot     string
	Miner            string
	Difficulty       string
	TotalDifficulty  string
	ExtraData        string
	Size             string
	GasLimit         string
	GasUsed          string
	Timestamp        string
	Transactions     []Transaction
	Uncles           []string
}

type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	ID      int             `json:"id"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func CallJSONRPC(method string, params []interface{}) (json.RawMessage, error) {
	ethereumConfig := config.GetConfig().Ethereum

	if ethereumConfig.Url == "" {
		return nil, errors.New("ethereum url is empty")
	}

	if method == "" {
		return nil, errors.New("method is empty")
	}

	reqBody, err := json.Marshal(JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	})

	if err != nil {
		return nil, errors.New("error marshalling request body, " + err.Error())
	}

	resp, err := http.Post(ethereumConfig.Url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, errors.New("error sending request, " + err.Error())
	}
	defer resp.Body.Close()

	var rpcResp JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return nil, errors.New("error decoding response body, " + err.Error())
	}

	if rpcResp.Error.Code != 0 {
		return nil, errors.New(rpcResp.Error.Message)
	}

	return rpcResp.Result, nil
}

func GetBlockNumber() (int, error) {
	result, err := CallJSONRPC("eth_blockNumber", []interface{}{})
	if err != nil {
		return 0, errors.New("error getting block number, " + err.Error())
	}

	var hexBlockNumber string
	if err := json.Unmarshal(result, &hexBlockNumber); err != nil {
		return 0, errors.New("error unmarshalling block number, " + err.Error())
	}

	blockNumber, err := strconv.ParseInt(hexBlockNumber[2:], 16, 64)
	if err != nil {
		return 0, errors.New("error parsing block number, " + err.Error())
	}

	return int(blockNumber), nil
}

func GetBlockByNumber(blockNumber int) (*Block, error) {
	result, err := CallJSONRPC("eth_getBlockByNumber", []interface{}{
		"0x" + strconv.FormatInt(int64(blockNumber), 16), true},
	)
	if err != nil {
		return nil, errors.New("error getting block, " + err.Error())
	}

	var block Block
	if err := json.Unmarshal(result, &block); err != nil {
		return nil, errors.New("error unmarshalling block, " + err.Error())
	}

	return &block, nil
}
