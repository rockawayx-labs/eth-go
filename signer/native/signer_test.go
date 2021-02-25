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

var b1e18, _ = new(big.Int).SetString("1000000000000000000", 0)
var b20e9, _ = new(big.Int).SetString("20000000000", 0)

type trx struct {
	nonce    uint64
	gasPrice *big.Int
	gasLimit uint64
	to       []byte
	value    *big.Int
	input    []byte
	chainID  uint64
}

func TestSigner_Hash(t *testing.T) {
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
			trx{9, b20e9, 21000, eth.MustNewAddress("0x3535353535353535353535353535353535353535"), b1e18, nil, 1},
			"4646464646464646464646464646464646464646464646464646464646464646",
			"37",
			"18515461264373351373200002665853028612451056578545711640558177340181847433846",
			"46948507304638947509940763649030358759909902576025900602547168820602576006531",
			nil,
		},

		{
			"parity odd",
			trx{9, b20e9, 21000, eth.MustNewAddress("0x3535353535353535353535353535353535353535"), b1e18, nil, 1},
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

			v, r, s, err := signer.Sign(test.in.nonce, test.in.to, test.in.value, test.in.gasLimit, test.in.gasPrice, test.in.input)

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

func TestSigner_ParityOdd(t *testing.T) {
	t.Skip("Used to generate a signature with parity being odd")

	for {
		privKey, err := eth.NewRandomPrivateKey()
		require.NoError(t, err)

		rawTx := trx{9, b20e9, 21000, eth.MustNewAddress("0x3535353535353535353535353535353535353535"), b1e18, nil, 1}

		data, err := rlp.Encode([]interface{}{
			rawTx.nonce,
			rawTx.gasPrice,
			rawTx.gasLimit,
			rawTx.to,
			rawTx.value,
			rawTx.input,
			big.NewInt(int64(rawTx.chainID)),
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
