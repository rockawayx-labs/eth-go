package eth

import (
	"math/big"
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

func TestMethod_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		cnt    string
		method *Method
	}{
		{
			name: "method with address",
			cnt:  `{"signature": "method(address)", "inputs": [{"type": "address", "value": "0xf8ac81dd843b9aca008302f8ac81dd843b9aca00"}]}`,
			method: &Method{
				Signature: "method(address)",
				Inputs: []*Input{
					{Type: "address", Value: MustNewAddress("f8ac81dd843b9aca008302f8ac81dd843b9aca00")},
				},
			},
		},
		{
			name: "method with uint112",
			cnt:  `{"signature": "method(uint112)", "inputs": [{"type": "uint112", "value": "10"}]}`,
			method: &Method{
				Signature: "method(uint112)",
				Inputs: []*Input{
					{Type: "uint112", Value: big.NewInt(10)},
				},
			},
		},
		{
			name: "method with uint256",
			cnt:  `{"signature": "method(uint256)", "inputs": [{"type": "uint256", "value": "23"}]}`,
			method: &Method{
				Signature: "method(uint256)",
				Inputs: []*Input{
					{Type: "uint256", Value: big.NewInt(23)},
				},
			},
		},
		{
			name: "method with string",
			cnt:  `{"signature": "method(string)", "inputs": [{"type": "string", "value": "hello world"}]}`,
			method: &Method{
				Signature: "method(string)",
				Inputs: []*Input{
					{Type: "string", Value: "hello world"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := NewMethodFromJSON([]byte(test.cnt))
			require.NoError(t, err)
			assert.Equal(t, test.method, m)
		})
	}
}
