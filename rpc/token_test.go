package rpc

import (
	"testing"

	"github.com/dfuse-io/eth-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTokenInfo(t *testing.T) {
	t.Skipf("Skipping token get info")
	client := NewClient("http://localhost:8545")
	addr := eth.MustNewAddress("0xd1c24bcabab5f01bcba2b44b6097ba506be67d6d")
	token, err := client.GetTokenInfo(addr)
	require.NoError(t, err)
	assert.Equal(t, &eth.Token{
		Name:               "Dogelon Mars",
		Symbol:             "ELON",
		Address:            addr,
		Decimals:           0,
		TotalSupply:        nil,
		IsEmptyName:        false,
		IsEmptyDecimal:     false,
		IsEmptySymbol:      false,
		IsEmptyTotalSupply: false,
	}, token)
}
