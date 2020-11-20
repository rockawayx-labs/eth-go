package rpc

import (
	"math/big"
	"testing"

	"github.com/dfuse-io/eth-go"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
)

func TestMarshalJSONRPC(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expected    string
		expectedErr error
	}{
		{"int", int(1), `"0x1"`, nil},
		{"uint", uint(1), `"0x1"`, nil},
		{"big.Int", *big.NewInt(1), `"0x1"`, nil},
		{"*big.Int", big.NewInt(1), `"0x1"`, nil},

		{"*eth.MethodCall", eth.MustNewMethodDef("name()").NewCall(), `"0x06fdde03"`, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := MarshalJSONRPC(test.in)
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.JSONEq(t, test.expected, string(actual))
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}
