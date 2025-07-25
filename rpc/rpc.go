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
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/streamingfast/eth-go"
	"github.com/streamingfast/logging"
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
	cache      Cache
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

func WithCache(cache Cache) Option {
	return func(client *Client) {
		client.cache = cache
	}
}

type LogsParams struct {
	// FromBlock is either block number encoded as a hexadecimal or tagged value which is one of
	// "latest" (`rpc.LatestBlock()`), "pending" (`rpc.PendingBlock()`) or "earliest" tags (`rpc.EarliestBlockRef`) (optional).
	FromBlock *BlockRef `json:"fromBlock,omitempty"`

	// ToBlock is either block number encoded as a hexadecimal or tagged value which is one of
	// "latest" (`LatestBlockRef()`), "pending" (`PendingBlockRef()`) or "earliest" tags (`EarliestBlockRef()`) (optional).
	ToBlock *BlockRef `json:"toBlock,omitempty"`

	// Address is the contract address or a list of addresses from which logs should originate (optional).
	Address eth.Address `json:"address,omitempty"`

	// Topics are order-dependent, each topic can also be an array of DATA with "or" options (optional).
	Topics *TopicFilter `json:"topics,omitempty"`
}

type TopicFilter struct {
	topics []TopicFilterExpr
}

func (f *TopicFilter) String() string {
	var elements []string
	for _, topic := range f.topics {
		elements = append(elements, topic.String())
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

func (f *TopicFilter) Append(in interface{}) {
	f.topics = append(f.topics, newTopicExpr(in))
}

func (f *TopicFilter) MarshalJSONRPC() ([]byte, error) {
	return MarshalJSONRPC(f.topics)
}

type TopicFilterExpr struct {
	exact *eth.Topic
	oneOf []eth.Topic
}

func (f TopicFilterExpr) String() string {
	if f.oneOf == nil {
		if f.exact == nil {
			return "null"
		}

		bytes := *f.exact
		return eth.Hex(bytes[:]).String()
	}

	var elements []string
	for _, topic := range f.oneOf {
		elements = append(elements, eth.Hex(topic[:]).String())
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

func (f TopicFilterExpr) MarshalJSONRPC() ([]byte, error) {
	if f.oneOf == nil {
		return MarshalJSONRPC(f.exact)
	}

	return MarshalJSONRPC(f.oneOf)
}

func NewTopicFilter(exprs ...interface{}) (out *TopicFilter) {
	if len(exprs) == 0 {
		return nil
	}

	topics := make([]TopicFilterExpr, len(exprs))
	for i, expr := range exprs {
		topics[i] = newTopicExpr(expr)
	}

	return &TopicFilter{topics: topics}
}

func newTopicExpr(expr interface{}) (out TopicFilterExpr) {
	switch v := expr.(type) {
	case TopicFilterExpr:
		return v
	case eth.Topic:
		return TopicFilterExpr{exact: &v}
	default:
		return ExactTopic(v)
	}
}

func ExactTopic(in interface{}) TopicFilterExpr {
	return TopicFilterExpr{exact: eth.LogTopic(in)}
}

func AnyTopic() TopicFilterExpr {
	return TopicFilterExpr{exact: nil}
}

func OneOfTopic(topics ...interface{}) (out TopicFilterExpr) {
	if len(topics) == 0 {
		panic("there must be at least one element to create a one of topic element")
	}

	out.oneOf = make([]eth.Topic, len(topics))
	for i, topic := range topics {
		logTopic := eth.LogTopic(topic)
		if logTopic == nil {
			panic("it's invalid to use nil value when building a one of topic element")
		}

		out.oneOf[i] = *logTopic
	}

	return
}

func (c *Client) Logs(ctx context.Context, params LogsParams) ([]*LogEntry, error) {
	result, err := c.DoRequest(ctx, "eth_getLogs", []interface{}{params})
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	var logs []*LogEntry
	err = json.Unmarshal([]byte(result), &logs)
	if err != nil {
		return nil, fmt.Errorf("unable to decode logs as JSON: %w", err)
	}

	return logs, nil
}

func (c *Client) GetBlockByNumber(ctx context.Context, blockNum uint64) (*Block, error) {
	resp, err := c.DoRequest(ctx, "eth_getBlockByNumber", []interface{}{eth.Uint64(blockNum), false})
	if err != nil {
		return nil, fmt.Errorf("unable to perform eth_getBlockByNumber request: %w", err)
	}

	var block *Block
	err = json.Unmarshal([]byte(resp), &block)
	if err != nil {
		return nil, fmt.Errorf("unable to decode block from JSON: %w", err)
	}

	return block, nil
}

func (c *Client) LatestBlockNum(ctx context.Context) (uint64, error) {
	resp, err := c.DoRequest(ctx, "eth_blockNumber", []interface{}{})
	if err != nil {
		return 0, fmt.Errorf("unable to perform eth_blockNumber request: %w", err)
	}

	value, err := strconv.ParseUint(resp, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse block number %s: %w", resp, err)
	}

	return value, nil
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

func (c *Client) Call(ctx context.Context, params CallParams) (string, error) {
	return c.callAtBlock(ctx, "eth_call", params, LatestBlock)
}

func (c *Client) CallAtBlock(ctx context.Context, params CallParams, blockAt *BlockRef) (string, error) {
	return c.callAtBlock(ctx, "eth_call", params, blockAt)
}

func (c *Client) EstimateGas(ctx context.Context, params CallParams) (string, error) {
	return c.callAtBlock(ctx, "eth_estimateGas", params, LatestBlock)
}

func (c *Client) callAtBlock(ctx context.Context, method string, params interface{}, blockAt *BlockRef) (string, error) {
	return c.DoRequest(ctx, method, []interface{}{params, blockAt})
}

// Depreacted: Use SendRawTransaction instead
func (c *Client) SendRaw(ctx context.Context, rawData []byte) (string, error) {
	return c.DoRequest(ctx, "eth_sendRawTransaction", []interface{}{rawData})
}

func (c *Client) SendRawTransaction(ctx context.Context, rawData []byte) (string, error) {
	return c.DoRequest(ctx, "eth_sendRawTransaction", []interface{}{rawData})
}

func (c *Client) ChainID(ctx context.Context) (*big.Int, error) {
	if c.chainID != nil {
		return c.chainID, nil
	}

	resp, err := c.DoRequest(ctx, "eth_chainId", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("unable to perform eth_chainId request: %w", err)
	}

	i := &big.Int{}
	_, ok := i.SetString(resp, 0)
	if !ok {
		return nil, fmt.Errorf("unable to parse chain id %s: %w", resp, err)
	}
	c.chainID = i
	return c.chainID, nil
}

func (c *Client) ProtocolVersion(ctx context.Context) (string, error) {
	resp, err := c.DoRequest(ctx, "eth_protocolVersion", []interface{}{})
	if err != nil {
		return "", fmt.Errorf("unable to perform eth_protocolVersion request: %w", err)
	}

	return resp, nil
}

type SyncingResp struct {
	StartingBlockNum eth.Uint64 `json:"starting_block_num"`
	CurrentBlockNum  eth.Uint64 `json:"current_block_num"`
	HighestBlockNum  eth.Uint64 `json:"highest_block_num"`
}

func (c *Client) Syncing(ctx context.Context) (*SyncingResp, error) {
	resp, err := c.DoRequest(ctx, "eth_syncing", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("unable to perform eth_syncing request: %w", err)
	}

	if resp == "false" {
		return nil, ErrFalseResp
	}

	out := &SyncingResp{}
	if err := json.Unmarshal([]byte(resp), out); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return out, nil
}

// TransactionReceipt fetches the receipt associated with the transaction's hash received. If the
// transaction is not found by the queried node, `nil, nil` is returned. If it's found, the receipt
// is decoded and `receipt, nil` is returned. Otherwise, the RPC error is returned if something went wrong.
func (c *Client) TransactionReceipt(ctx context.Context, hash eth.Hash) (out *TransactionReceipt, err error) {
	resp, err := c.DoRequest(ctx, "eth_getTransactionReceipt", []interface{}{hash})
	if err != nil {
		return nil, fmt.Errorf("unable to perform eth_getTransactionCount request: %w", err)
	}

	if resp == "" {
		return nil, nil
	}

	err = json.Unmarshal([]byte(resp), &out)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return out, nil
}

func (c *Client) GetTransactionCount(ctx context.Context, accountAddr eth.Address) (uint64, error) {
	return c.Nonce(ctx, accountAddr)
}

func (c *Client) Nonce(ctx context.Context, accountAddr eth.Address) (uint64, error) {
	resp, err := c.DoRequest(ctx, "eth_getTransactionCount", []interface{}{accountAddr.Pretty(), LatestBlock})
	if err != nil {
		return 0, fmt.Errorf("unable to perform eth_getTransactionCount request: %w", err)
	}

	nonce, err := strconv.ParseUint(strings.TrimPrefix(resp, "0x"), 16, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse nonce %s: %w", resp, err)
	}
	return nonce, nil
}

func (c *Client) GetBalance(ctx context.Context, accountAddr eth.Address) (*eth.TokenAmount, error) {
	resp, err := c.DoRequest(ctx, "eth_getBalance", []interface{}{accountAddr.Pretty(), LatestBlock})
	if err != nil {
		return nil, fmt.Errorf("unable to perform eth_getBalance request: %w", err)
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

func (c *Client) GasPrice(ctx context.Context) (*big.Int, error) {
	resp, err := c.DoRequest(ctx, "eth_gasPrice", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("unable to perform eth_gasPrice request: %w", err)
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
	Method  string        `json:"method"`
	decoder ResponseDecoder

	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

type RPCResponse struct {
	Content string
	ID      int
	Err     error
	decoder ResponseDecoder
}

func (res *RPCResponse) CopyDecoder(req *RPCRequest) {
	res.decoder = req.decoder
}

func (res *RPCResponse) Empty() bool {
	return res.Content == "0x"
}

func (res *RPCResponse) Deterministic() bool {
	if res.Err == nil {
		return true
	}
	if rpcErr, ok := res.Err.(*ErrResponse); ok {
		return IsDeterministicError(rpcErr)
	}
	return false
}

func (res *RPCResponse) Decode() ([]interface{}, error) {
	if res.decoder == nil {
		return nil, fmt.Errorf("no decoder in RPC response")
	}
	if res.Err != nil {
		return nil, fmt.Errorf("error in response, cannot decode")
	}
	if res.Empty() {
		return nil, fmt.Errorf("empty response, cannot decode")
	}
	return res.decoder(eth.MustNewHex(res.Content))
}

type ETHCallOption func(*ETHCall)

func AtBlockNum(num uint64) ETHCallOption {
	return func(c *ETHCall) {
		c.atExpr = num
	}
}

func NewETHCall(to eth.Address, methodDef *eth.MethodDef, options ...ETHCallOption) *ETHCall {
	c := &ETHCall{
		params: CallParams{
			To:   to,
			Data: methodDef.NewCall().MustEncode(),
		},
		methodDef:       methodDef,
		responseDecoder: methodDef.DecodeOutput,
		atExpr:          LatestBlock,
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

type ETHCall struct {
	params          CallParams
	methodDef       *eth.MethodDef
	atExpr          interface{}
	responseDecoder ResponseDecoder
}

func (c *ETHCall) ToRequest() *RPCRequest {
	return &RPCRequest{
		Params:  []interface{}{c.params, c.atExpr},
		decoder: c.responseDecoder,
		Method:  "eth_call",
	}
}

type ResponseDecoder func([]byte) ([]interface{}, error)

func (c *Client) DoRequests(ctx context.Context, reqs []*RPCRequest) ([]*RPCResponse, error) {
	logger := logging.Logger(ctx, zlog).With(zap.Strings("methods", methodsFromRPCRequests(reqs)))

	// sanitize reqs
	var lastID int
	// we need IDs to be sorted
	for _, req := range reqs {
		lastID++
		req.ID = lastID
		req.JSONRPC = "2.0"
	}

	reqsBytes, err := MarshalJSONRPC(&reqs)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal json_rpc requests: %w", err)
	}
	if tracer.Enabled() {
		logger.Debug("json_rpc requests", zap.Stringer("requests", eth.Hex(reqsBytes)))
	}

	var cacheKey string
	if c.cache != nil {
		cacheKey = hex.EncodeToString(sha256.New().Sum(reqsBytes))
	}

	resp, err := c.doRequest(ctx, logger, reqsBytes, cacheKey)
	if err != nil {
		return nil, err
	}

	results, err := parseRPCResults(logger, resp)
	if err != nil {
		return nil, err
	}

	if len(results) != len(reqs) {
		logger.Warn("invalid number of results", zap.Int("len_results", len(results)), zap.Int("len_reqs", len(reqs)))
		return nil, fmt.Errorf("invalid number of results")
	}

	if c.cache != nil {
		deterministic := true
		for i, req := range reqs {
			results[i].decoder = req.decoder
			if !results[i].Deterministic() {
				deterministic = false
				break
			}
		}

		if deterministic {
			c.cache.Set(ctx, cacheKey, resp)
		}
	}

	return results, nil
}

func (c *Client) DoRequest(ctx context.Context, method string, params []interface{}) (string, error) {
	logger := logging.Logger(ctx, zlog).With(zap.String("method", method))

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

	if tracer.Enabled() {
		logger.Debug("json_rpc request", zap.String("request", string(reqCnt)))
	}

	cacheKey := hex.EncodeToString(sha256.New().Sum(reqCnt))

	resp, err := c.doRequest(ctx, logger, reqCnt, cacheKey)
	if err != nil {
		return "", err
	}

	results, err := parseRPCResults(logger, resp)
	if err != nil {
		return "", err
	}
	if len(results) != 1 {
		return "", fmt.Errorf("received no result than number of requests")
	}

	if c.cache != nil && results[0].Deterministic() {
		c.cache.Set(ctx, cacheKey, resp)
	}

	return results[0].Content, results[0].Err
}

func (c *Client) doRequest(ctx context.Context, logger *zap.Logger, reqsBytes []byte, cacheKey string) ([]byte, error) {
	if c.cache != nil {
		cachedData, found := c.cache.Get(ctx, cacheKey)
		if found {
			if tracer.Enabled() {
				logger.Debug("retrieve request's response from cache", zap.String("key", cacheKey))
			}

			return cachedData, nil
		}
	}

	body := bytes.NewBuffer(reqsBytes)

	resp, err := c.post(ctx, c.URL, body)
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

	if tracer.Enabled() {
		logger.Debug("json_rpc call response", zap.String("response_body", string(bodyBytes)))
	}

	return bodyBytes, nil
}

func (c *Client) post(ctx context.Context, url string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(req)
}

func parseRPCResults(logger *zap.Logger, in []byte) ([]*RPCResponse, error) {
	responses := []gjson.Result{}

	// we may receive `[{resp1},{resp2}]` OR `{resp}`
	parsed := gjson.ParseBytes(in)
	if parsed.IsArray() {
		responses = parsed.Array()
	} else {
		responses = append(responses, parsed)
	}

	var out []*RPCResponse
	for _, response := range responses {
		rpcErrorResult := response.Get("error")
		if !rpcErrorResult.Exists() {
			out = append(out, &RPCResponse{Content: response.Get("result").String(), ID: int(hex2uint64(response.Get("id").String()))})
			continue
		}

		content := rpcErrorResult.Raw
		if tracer.Enabled() {
			logger.Error("json_rpc call response error",
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

		out = append(out, &RPCResponse{Err: rpcErr, ID: int(hex2uint64(response.Get("id").String()))})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})

	return out, nil
}

func methodsFromRPCRequests(requests []*RPCRequest) (out []string) {
	out = make([]string, len(requests))
	for i, v := range requests {
		out[i] = v.Method
	}

	return
}

func hex2uint64(hexStr string) uint64 {
	cleaned := strings.Replace(hexStr, "0x", "", -1)
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return result
}
