package rpc

import (
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/dfuse-io/eth-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTokenInfo(t *testing.T) {
	t.Skip()
	client := NewClient("http://localhost:8545")
	addr := eth.MustNewAddress("0xd1c24bcabab5f01bcba2b44b6097ba506be67d6d")
	token, err := client.GetTokenInfo(addr)
	require.NoError(t, err)
	assert.Equal(t, &eth.Token{
		Name:               "Dogelon Mars",
		Symbol:             "ELON",
		Address:            addr,
		TotalSupply:        big.NewInt(1000000000000000),
		Decimals:           0,
		IsEmptyName:        false, //		IsEmptyDecimal:     false, //		IsEmptySymbol:      false,
		IsEmptyTotalSupply: false,
	}, token)

}

func hex2int(hexStr string) uint64 {
	cleaned := strings.Replace(hexStr, "0x", "", -1)
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return uint64(result)
}

func TestGetAtBlockNum(t *testing.T) {
	//t.Skip()
	client := NewClient("http://localhost:8545")
	addr := eth.MustNewAddress("0xd1c24bcabab5f01bcba2b44b6097ba506be67d6d")
	results, err := client.DoRequests([]RPCRequest{
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

	toInt := hex2int(results[0].content)
	assert.Greater(t, toInt, uint64(8000000))
	assert.Equal(t, RPCResult{content: "0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000c446f67656c6f6e204d6172730000000000000000000000000000000000000000"}, results[1])
}
