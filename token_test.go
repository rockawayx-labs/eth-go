// Copyright 2021 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eth

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
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
