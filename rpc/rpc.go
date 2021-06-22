// Copyright 2021 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/dfuse-io/eth-go"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

var ErrFalseResp = errors.New("false response")

type Option func(*Client)

// TODO: refactor to use mux rpc
type Client struct {
	URL     string
	chainID *big.Int

	httpClient *http.Client
}

func NewClient(url string, opts ...Option) *Client {
	c := &Client{
		URL: url,
	}

	c.httpClient = http.DefaultClient

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

type CallParams struct {
	// From the address the transaction is sent from (optional).
	From eth.Address `json:"from,omitempty"`
	// To the address the transaction is directed to (required).
	To eth.Address `json:"to,omitempty"`
	// GasLimit Integer of the gas provided for the transaction execution. eth_call consumes zero gas, but this parameter may be needed by some executions (optional).
	GasLimit uint64 `json:"gas,omitempty"`
	// GasPrice big integer of the gasPrice used for each paid gas (optional).
	GasPrice *big.Int `json:"gasPrice,omitempty"`
	// Value big integer of the value sent with this transaction (optional).
	Value *big.Int `json:"value,omitempty"`
	// Hash of the method signature and encoded parameters or any object that implements `MarshalJSONRPC` and serialize to a byte array, for details see Ethereum Contract ABI in the Solidity documentation (optional).
	Data interface{} `json:"data,omitempty"`
}

func (c *Client) Call(params CallParams) (string, error) {
	return c.callAtBlock("eth_call", params, "latest")
}

func (c *Client) CallAtBlock(params CallParams, blockAt string) (string, error) {
	return c.callAtBlock("eth_call", params, blockAt)
}

func (c *Client) EstimateGas(params CallParams) (string, error) {
	return c.callAtBlock("eth_estimateGas", params, "latest")
}

func (c *Client) callAtBlock(method string, params interface{}, blockAt string) (string, error) {
	return c.DoRequest(method, []interface{}{params, blockAt})
}

func (c *Client) SendRaw(rawData []byte) (string, error) {
	return c.DoRequest("eth_sendRawTransaction", []interface{}{rawData})
}

func (c *Client) ChainID() (*big.Int, error) {
	if c.chainID != nil {
		return c.chainID, nil
	}

	resp, err := c.DoRequest("eth_chainId", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("unale to perform eth_chainId request: %w", err)
	}

	i := &big.Int{}
	_, ok := i.SetString(resp, 0)
	if !ok {
		return nil, fmt.Errorf("unable to parse chain id %s: %w", resp, err)
	}
	c.chainID = i
	return c.chainID, nil
}

func (c *Client) ProtocolVersion() (string, error) {
	resp, err := c.DoRequest("eth_protocolVersion", []interface{}{})
	if err != nil {
		return "", fmt.Errorf("unale to perform eth_protocolVersion request: %w", err)
	}

	return resp, nil
}

type SyncingResp struct {
	StartingBlockNum uint64 `json:"starting_block_num"`
	CurrentBlockNum  uint64 `json:"current_block_num"`
	HighestBlockNum  uint64 `json:"highest_block_num"`
}

func (c *Client) Syncing() (*SyncingResp, error) {
	resp, err := c.DoRequest("eth_syncing", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("unale to perform eth_syncing request: %w", err)
	}

	if resp == "false" {
		return nil, ErrFalseResp
	}
	out := &SyncingResp{}

	out.StartingBlockNum, err = strconv.ParseUint(strings.TrimPrefix(gjson.GetBytes([]byte(resp), "startingBlock").String(), "0x"), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse starting block num %s: %w", resp, err)
	}

	out.CurrentBlockNum, err = strconv.ParseUint(strings.TrimPrefix(gjson.GetBytes([]byte(resp), "currentBlock").String(), "0x"), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse current block num %s: %w", resp, err)
	}

	out.HighestBlockNum, err = strconv.ParseUint(strings.TrimPrefix(gjson.GetBytes([]byte(resp), "highestBlock").String(), "0x"), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse current block num %s: %w", resp, err)
	}

	return out, nil
}

func (c *Client) Nonce(accountAddr eth.Address) (uint64, error) {
	resp, err := c.DoRequest("eth_getTransactionCount", []interface{}{accountAddr.Pretty(), "latest"})
	if err != nil {
		return 0, fmt.Errorf("unale to perform eth_getTransactionCount request: %w", err)
	}

	nonce, err := strconv.ParseUint(strings.TrimPrefix(resp, "0x"), 16, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse nonce %s: %w", resp, err)
	}
	return nonce, nil

}

func (c *Client) GetBalance(accountAddr eth.Address) (*eth.TokenAmount, error) {
	resp, err := c.DoRequest("eth_getBalance", []interface{}{accountAddr.Pretty(), "latest"})
	if err != nil {
		return nil, fmt.Errorf("unale to perform eth_getBalance request: %w", err)
	}

	v, ok := new(big.Int).SetString(strings.TrimPrefix(resp, "0x"), 16)
	if !ok {
		return nil, fmt.Errorf("unable to parse balance %s: %w", resp, err)
	}

	return &eth.TokenAmount{
		Amount: v,
		Token:  eth.ETHToken,
	}, nil
}

func (c *Client) GasPrice() (*big.Int, error) {
	resp, err := c.DoRequest("eth_gasPrice", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("unale to perform eth_gasPrice request: %w", err)
	}

	i := &big.Int{}
	_, ok := i.SetString(resp, 0)
	if !ok {
		return nil, fmt.Errorf("unable to parse gas price %s: %w", resp, err)
	}

	return i, nil
}

type RPCRequest struct {
	Params  []interface{} `json:"params"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	ID      int           `json:"id"`
}

type RPCResult struct {
	content string
	err     error
}

func (c *Client) DoRequests(reqs []RPCRequest) ([]RPCResult, error) {
	// sanitize reqs
	var lastID int
	for _, req := range reqs {
		if req.ID == 0 {
			lastID++
			req.ID = lastID
		}
		if req.JSONRPC == "" {
			req.JSONRPC = "2.0"
		}
	}

	reqsBytes, err := MarshalJSONRPC(&reqs)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal json_rpc requests: %w", err)
	}

	if traceEnabled {
		zlog.Debug("json_rpc requests", zap.String("requests", string(reqsBytes)))
	}

	resp, err := c.doRequest(bytes.NewBuffer(reqsBytes))
	if err != nil {
		return nil, err
	}

	return parseRPCResults(resp)
}

func (c *Client) DoRequest(method string, params []interface{}) (string, error) {
	req := RPCRequest{
		Params:  params,
		JSONRPC: "2.0",
		Method:  method,
		ID:      1,
	}
	reqCnt, err := MarshalJSONRPC(&req)
	if err != nil {
		return "", fmt.Errorf("unable to marshal json_rpc request: %w", err)
	}

	if traceEnabled {
		zlog.Debug("json_rpc request", zap.String("request", string(reqCnt)))
	}

	resp, err := c.doRequest(bytes.NewBuffer(reqCnt))
	if err != nil {
		return "", err
	}

	results, err := parseRPCResults(resp)
	if err != nil {
		return "", err
	}
	return results[0].content, results[0].err
}

func (c *Client) doRequest(body *bytes.Buffer) ([]byte, error) {
	resp, err := c.httpClient.Post(c.URL, "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("sending request to json_rpc endpoint: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error in response: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read json_rpc response body: %w", err)
	}

	if traceEnabled {
		zlog.Debug("json_rpc call response", zap.String("response_body", string(bodyBytes)))
	}
	return bodyBytes, nil
}

func parseRPCResults(in []byte) ([]RPCResult, error) {
	responses := []gjson.Result{}
	parsed := gjson.ParseBytes(in)
	if parsed.IsArray() {
		responses = parsed.Array()
	} else {
		responses = append(responses, parsed)
	}

	var out []RPCResult
	for _, response := range responses {
		rpcErrorResult := response.Get("error")
		if !rpcErrorResult.Exists() {
			out = append(out, RPCResult{content: response.Get("result").String()})
			continue
		}

		content := rpcErrorResult.Raw
		if traceEnabled {
			zlog.Error("json_rpc call response error",
				zap.String("response_body", string(in)),
				zap.String("error", content),
			)
		}

		rpcErr := &ErrResponse{}
		err := json.Unmarshal([]byte(content), rpcErr)
		if err != nil {
			// We were not able to deserialize to RPC error, return the whole thing with an error
			return nil, fmt.Errorf("json_rpc returned error: %s", rpcErrorResult)
		}

		out = append(out, RPCResult{err: rpcErr})
	}
	return out, nil
}
