package eth

import (
	"encoding/hex"
	"github.com/dfuse-io/logging"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func init() {
	logging.TestingOverride()
}

func bs(in string) *big.Int {
	v := &big.Int{}
	v.SetString(in, 10)
	return v
}

func bigToUint(in uint64) *big.Int {
	z1 := &big.Int{}
	z1.SetUint64(in) // z1 := 123
	return z1
}

func stringToByte(t *testing.T, address string) []byte {
	out, err := hex.DecodeString(address)
	require.NoError(t, err)

	return out
}
