package rpc

import (
	"bytes"
	"encoding/hex"
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

// TODO: refactor to use mux rpc
type Client struct {
	URL     string
	chainID *big.Int
}

func NewClient(url string) *Client {
	return &Client{
		URL: url,
	}
}

type ethCallData struct {
	Params  []interface{} `json:"params"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	ID      int           `json:"id"`
}

/*
The trxObject object:
from: (optional) The address the transaction is sent from.
to: The address the transaction is directed to.
gas: (optional) Integer of the gas provided for the transaction execution. eth_call consumes zero gas, but this parameter may be needed by some executions.
gasPrice: (optional) Integer of the gasPrice used for each paid gas
value:	(optional) Integer of the value sent with this transaction
data:	(optional) Hash of the method signature and encoded parameters. For details see Ethereum Contract ABI in the Solidity documentation
*/
func (c *Client) Call(trxObject map[string]string) (string, error) {
	return c.CallAtBlock(trxObject, "latest")
}

func (c *Client) CallAtBlock(trxObject map[string]string, blockAt string) (string, error) {
	return c.rpcCall("eth_call", []interface{}{trxObject, blockAt})
}

func (c *Client) SendRaw(rawData []byte) (string, error) {
	return c.rpcCall("eth_sendRawTransaction", []interface{}{eth.PrefixedHex(hex.EncodeToString(rawData))})
}

func (c *Client) ChainID() (*big.Int, error) {
	if c.chainID != nil {
		return c.chainID, nil
	}

	resp, err := c.rpcCall("eth_chainId", []interface{}{})
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
	resp, err := c.rpcCall("eth_protocolVersion", []interface{}{})
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
	resp, err := c.rpcCall("eth_syncing", []interface{}{})
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
	resp, err := c.rpcCall("eth_getTransactionCount", []interface{}{accountAddr.Pretty(), "latest"})
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
	resp, err := c.rpcCall("eth_getBalance", []interface{}{accountAddr.Pretty(), "latest"})
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
	resp, err := c.rpcCall("eth_gasPrice", []interface{}{})
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

type rpcRequest struct {
	Params  []interface{} `json:"params"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	ID      int           `json:"id"`
}

func (c *Client) rpcCall(method string, params []interface{}) (string, error) {
	req := rpcRequest{
		Params:  params,
		JSONRPC: "2.0",
		Method:  method,
		ID:      1,
	}
	reqCnt, err := json.Marshal(&req)
	if err != nil {
		return "", fmt.Errorf("unable to marshall json_rpc request: %w", err)
	}

	if traceEnabled {
		zlog.Debug("json_rpc request",
			zap.String("request", string(reqCnt)),
		)
	}

	return c.doRequest(bytes.NewBuffer(reqCnt))
}

func (c *Client) doRequest(body *bytes.Buffer) (string, error) {
	resp, err := http.Post(c.URL, "application/json", body)
	if err != nil {
		return "", fmt.Errorf("sending request to json_rpc endpoint: %w", err)
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("error in response: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read json_rpc response body: %w", err)
	}

	if traceEnabled {
		zlog.Debug("json_rpc call response", zap.String("response_body", string(bodyBytes)))
	}

	error := gjson.GetBytes(bodyBytes, "error").String()
	if error != "" {
		if traceEnabled {
			zlog.Error("json_rpc call response error",
				zap.String("response_body", string(bodyBytes)),
				zap.String("error", error),
			)
		}
		return "", fmt.Errorf("json_rpc returned error: %s", error)
	}

	result := gjson.GetBytes(bodyBytes, "result").String()
	return result, nil

}
