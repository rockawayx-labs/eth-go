package eth

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeccak256(t *testing.T) {
	tests := []struct {
		in        string
		expectOut string
	}{
		{
			in:        "Pregnant(address,uint256,uint256,uint256)",
			expectOut: "241ea03ca20251805084d27d4440371c34a0b85ff108f6bb5611248f73818b80",
		},
		{
			in:        "Transfer(address,address,uint256)",
			expectOut: "ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		},
		{
			in:        "Approval(address,address,uint256)",
			expectOut: "8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925",
		},
		{
			in:        "Birth(address,uint256,uint256,uint256,uint256)",
			expectOut: "0a5311bd2a6608f08a180df2ee7c5946819a649b204b554bb8e39825b2c50ad5",
		},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			out := Keccak256([]byte(test.in))
			assert.Equal(t, test.expectOut, hex.EncodeToString(out))
		})
	}
}
