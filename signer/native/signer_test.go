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

package native

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/dfuse-io/eth-go"
	"github.com/dfuse-io/eth-go/rlp"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

var b1 = big.NewInt(1)
var b1e18, _ = new(big.Int).SetString("1000000000000000000", 0)
var b20e9, _ = new(big.Int).SetString("20000000000", 0)

type trx struct {
	nonce    uint64
	gasPrice *big.Int
	gasLimit uint64
	to       []byte
	value    *big.Int
	input    []byte
	chainID  *big.Int
}

func TestSigner_Signature(t *testing.T) {
	tests := []struct {
		name        string
		in          trx
		privateKey  string
		expectedV   string
		expectedR   string
		expectedS   string
		expectedErr error
	}{
		{
			"parity even",
			trx{9, b20e9, 21000, eth.MustNewAddress("0x3535353535353535353535353535353535353535"), b1e18, nil, b1},
			"4646464646464646464646464646464646464646464646464646464646464646",
			"37",
			"18515461264373351373200002665853028612451056578545711640558177340181847433846",
			"46948507304638947509940763649030358759909902576025900602547168820602576006531",
			nil,
		},

		{
			"parity odd",
			trx{9, b20e9, 21000, eth.MustNewAddress("0x3535353535353535353535353535353535353535"), b1e18, nil, b1},
			"5141949acb9d2c7f5667c25428d61d290ced0d721eeb026c291da8018c0e8f50",
			"38",
			"23672060004711451020640132671354144284260356544239847582982888380471096525220",
			"42942693600234235375806513579982379680684107821147881578372863492373724509763",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			priv, err := eth.NewPrivateKey(test.privateKey)
			require.NoError(t, err)

			signer, err := NewPrivateKeySigner(zlog, test.in.chainID, priv)
			require.NoError(t, err)

			v, r, s, err := signer.Signature(test.in.nonce, test.in.to, test.in.value, test.in.gasLimit, test.in.gasPrice, test.in.input)

			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, bigString(t, test.expectedV), v, "Signature V")
				assert.Equal(t, bigString(t, test.expectedR), r, "Signature R")
				assert.Equal(t, bigString(t, test.expectedS), s, "Signature S")
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func TestSigner_Sign(t *testing.T) {
	tests := []struct {
		name        string
		in          trx
		privateKey  string
		expected    string
		expectedErr error
	}{
		{
			"parity even",
			trx{0, bigString(t, "0x3b9aca00"), 21000, eth.MustNewAddress("0x17A98d2b11Dfb784e63337d2170e21cf5DD04631"), bigString(t, "0x16345785d8a0000"), nil, b1},
			"616e6769652e6a6a706572657a616775696e6167612e6574682e6c696e6b0d0a",
			"f86b80843b9aca008252089417a98d2b11dfb784e63337d2170e21cf5dd0463188016345785d8a00008025a02e47aa4c37e7003af4d3b7d20265691b6c03baba509c0556d21acaca82876cb4a01b5711b8c801584c7875370ed2e9b60260b390cdb63cf57fa6d77899102279a0",
			nil,
		},

		// {
		// 	"parity odd",
		// 	trx{9, b20e9, 21000, eth.MustNewAddress("0x3535353535353535353535353535353535353535"), b1e18, nil, b1},
		// 	"5141949acb9d2c7f5667c25428d61d290ced0d721eeb026c291da8018c0e8f50",
		// 	"38",
		// 	nil,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			priv, err := eth.NewPrivateKey(test.privateKey)
			require.NoError(t, err)

			signer, err := NewPrivateKeySigner(zlog, test.in.chainID, priv)
			require.NoError(t, err)

			actual, err := signer.Sign(test.in.nonce, test.in.to, test.in.value, test.in.gasLimit, test.in.gasPrice, test.in.input)

			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expected, hex.EncodeToString(actual))
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func TestSigner_ParityOdd(t *testing.T) {
	t.Skip("Used to generate a signature with parity bit being odd")

	for {
		privKey, err := eth.NewRandomPrivateKey()
		require.NoError(t, err)

		rawTx := trx{9, b20e9, 21000, eth.MustNewAddress("0x3535353535353535353535353535353535353535"), b1e18, nil, b1}

		data, err := rlp.Encode([]interface{}{
			rawTx.nonce,
			rawTx.gasPrice,
			rawTx.gasLimit,
			rawTx.to,
			rawTx.value,
			rawTx.input,
			rawTx.chainID,
			uint64(0),
			uint64(0),
		})

		require.NoError(t, err)

		hash := eth.Keccak256(data)

		btcecPrivKey := (*btcec.PrivateKey)(privKey.ToECDSA())
		compressedSignature, err := btcec.SignCompact(btcec.S256(), btcecPrivKey, hash, false)
		require.NoError(t, err)

		if compressedSignature[0] != 27 {
			fmt.Printf("Parity value is odd (recovery id %d) for private key %s (signature %s)\n", compressedSignature[0], privKey.String(), hex.EncodeToString(compressedSignature))
			require.NoError(t, errors.New("parity odd"))
		}
	}
}

func bigString(t *testing.T, value string) *big.Int {
	t.Helper()

	out := new(big.Int)
	out, ok := out.SetString(value, 0)
	assert.True(t, ok, "In value %q is not a valid big.Int string", value)

	return out
}
