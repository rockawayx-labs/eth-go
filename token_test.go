package eth

import (
	"math/big"
	"testing"

	"gotest.tools/assert"
)

func TestTokenAmount_Format(t *testing.T) {
	tests := []struct {
		name       string
		token      *Token
		amount     *big.Int
		exectValue string
	}{
		{
			name: "positive whole amount",
			token: &Token{
				Name:     "Wrapped Ethereum",
				Symbol:   "WETH",
				Address:  MustNewAddress("0x22f4d9ba781bc8f97dc94ea93554c7acb9f8bd44"),
				Decimals: 5,
			},
			amount:     bigString(t, "1800000"),
			exectValue: "18.0000 WETH",
		},
		{
			name: "positive amount",
			token: &Token{
				Name:     "Wrapped Ethereum",
				Symbol:   "WETH",
				Address:  MustNewAddress("0x22f4d9ba781bc8f97dc94ea93554c7acb9f8bd44"),
				Decimals: 5,
			},
			amount:     bigString(t, "1826731"),
			exectValue: "18.2673 WETH",
		},
		{
			name: "negative amount",
			token: &Token{
				Name:     "Wrapped Ethereum",
				Symbol:   "WETH",
				Address:  MustNewAddress("0x22f4d9ba781bc8f97dc94ea93554c7acb9f8bd44"),
				Decimals: 5,
			},
			amount:     new(big.Int).Sub(big.NewInt(500000), big.NewInt(800000)),
			exectValue: "-3.0000 WETH",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenAmount := &TokenAmount{
				Amount: test.amount,
				Token:  test.token,
			}
			assert.Equal(t, test.exectValue, tokenAmount.String())
		})
	}
}
