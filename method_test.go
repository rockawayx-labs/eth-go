package eth

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMethodCall_AppendArgFromString(t *testing.T) {
	tests := []struct {
		name             string
		signature        string
		inputs           []string
		expectMethodDef  *MethodDef
		expectMethodCall *MethodCall
	}{
		{
			name:      "testing bytes",
			signature: "method(bytes)",
			inputs:    []string{"0xaabbcc"},
			expectMethodDef: &MethodDef{
				Name:       "method",
				Parameters: []*MethodParameter{{TypeName: "bytes"}},
			},
			expectMethodCall: &MethodCall{
				MethodDef: &MethodDef{Name: "method", Parameters: []*MethodParameter{{TypeName: "bytes"}}},
				Data: []interface{}{
					[]byte{0xaa, 0xbb, 0xcc},
				},
			},
		},
		{
			name:      "testing address",
			signature: "method(address)",
			inputs:    []string{"0x5a0b54d5dc17e0aadc383d2db43b0a0d3e029c4c"},
			expectMethodDef: &MethodDef{
				Name:       "method",
				Parameters: []*MethodParameter{{TypeName: "address"}},
			},
			expectMethodCall: &MethodCall{
				MethodDef: &MethodDef{Name: "method", Parameters: []*MethodParameter{{TypeName: "address"}}},
				Data: []interface{}{
					MustNewAddress("5a0b54d5dc17e0aadc383d2db43b0a0d3e029c4c"),
				},
			},
		},
		{
			name:      "testing uint64",
			signature: "method(uint64)",
			inputs:    []string{"13"},
			expectMethodDef: &MethodDef{
				Name:       "method",
				Parameters: []*MethodParameter{{TypeName: "uint64"}},
			},
			expectMethodCall: &MethodCall{
				MethodDef: &MethodDef{Name: "method", Parameters: []*MethodParameter{{TypeName: "uint64"}}},
				Data: []interface{}{
					uint64(13),
				},
			},
		},
		{
			name:      "testing uint112",
			signature: "method(uint112)",
			inputs:    []string{"123456789"},
			expectMethodDef: &MethodDef{
				Name:       "method",
				Parameters: []*MethodParameter{{TypeName: "uint112"}},
			},
			expectMethodCall: &MethodCall{
				MethodDef: &MethodDef{Name: "method", Parameters: []*MethodParameter{{TypeName: "uint112"}}},
				Data: []interface{}{
					big.NewInt(123456789),
				},
			},
		},
		{
			name:      "testing uint256",
			signature: "method(uint256)",
			inputs:    []string{"123456789"},
			expectMethodDef: &MethodDef{
				Name:       "method",
				Parameters: []*MethodParameter{{TypeName: "uint256"}},
			},
			expectMethodCall: &MethodCall{
				MethodDef: &MethodDef{Name: "method", Parameters: []*MethodParameter{{TypeName: "uint256"}}},
				Data: []interface{}{
					big.NewInt(123456789),
				},
			},
		},
		{
			name:      "testing bool",
			signature: "method(bool)",
			inputs:    []string{"true"},
			expectMethodDef: &MethodDef{
				Name:       "method",
				Parameters: []*MethodParameter{{TypeName: "bool"}},
			},
			expectMethodCall: &MethodCall{
				MethodDef: &MethodDef{Name: "method", Parameters: []*MethodParameter{{TypeName: "bool"}}},
				Data: []interface{}{
					true,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			methodDef, err := NewMethodDef(test.signature)
			require.NoError(t, err)
			assert.Equal(t, test.expectMethodDef, methodDef)

			methodCall := methodDef.NewCall()
			for _, input := range test.inputs {
				err := methodCall.AppendArgFromString(input)
				require.NoError(t, err)
			}

			assert.Equal(t, test.expectMethodCall, methodCall)
		})
	}
}

func TestMethodCall_AppendArgFromStringTooManyArgs(t *testing.T) {
	methodDef, err := NewMethodDef("method(address)")
	require.NoError(t, err)
	assert.Equal(t, &MethodDef{
		Name:       "method",
		Parameters: []*MethodParameter{{TypeName: "address"}},
	}, methodDef)

	methodCall := methodDef.NewCall()
	err = methodCall.AppendArgFromString("0x5a0b54d5dc17e0aadc383d2db43b0a0d3e029c4c")
	require.NoError(t, err)

	err = methodCall.AppendArgFromString("0x5a0b54d5dc17e0aadc383d2db43b0a0d3e029c4c")
	assert.Error(t, err)
}
