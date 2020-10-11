package eth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAddress(t *testing.T) {
	method, err := NewMethodFromSignature("swapExactTokensForTokens(uint256,uint256,address[],address,uint256)")
	require.NoError(t, err)
	assert.Equal(t, &Method{
		Signature: "swapExactTokensForTokens(uint256,uint256,address[],address,uint256)",
		Inputs: []*Input{
			{Type: "uint256"},
			{Type: "uint256"},
			{Type: "address[]"},
			{Type: "address"},
			{Type: "uint256"},
		},
	},
		method)
}
