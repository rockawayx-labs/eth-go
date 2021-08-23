package rpc

import (
	"math/big"
	"testing"

	"github.com/streamingfast/eth-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTokenInfo(t *testing.T) {
	t.Skip() //requires a full archive node at localhost:8545 on eth-mainnet
	client := NewClient("http://localhost:8545")
	addr := eth.MustNewAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")

	token, err := client.GetTokenInfo(addr, 5000000)
	require.NoError(t, err)
	assert.Equal(t, &eth.Token{
		Name:        "Tether USD",
		Symbol:      "USDT",
		Address:     addr,
		TotalSupply: big.NewInt(60109970000000),
		Decimals:    6,
	}, token)
}

func TestGetAtBlockNum(t *testing.T) {
	t.Skip()
	client := NewClient("http://localhost:8545")
	addr := eth.MustNewAddress("0xd1c24bcabab5f01bcba2b44b6097ba506be67d6d")
	results, err := client.DoRequests([]*RPCRequest{
		{
			Params: []interface{}{},
			Method: "eth_blockNumber",
		},
		{
			Params: []interface{}{CallParams{To: addr, Data: nameCallData}, "latest"},
			Method: "eth_call",
		},
	},
	)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	require.NoError(t, results[0].err)
	assert.Equal(t, 1, results[0].ID)
	assert.Equal(t, 2, results[1].ID)

	toInt := int(hex2uint64(results[0].content))
	assert.Greater(t, toInt, 8000000)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000c446f67656c6f6e204d6172730000000000000000000000000000000000000000", results[1].content)
}
