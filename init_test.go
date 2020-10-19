package eth

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/dfuse-io/logging"
	"github.com/stretchr/testify/require"
)

func init() {
	logging.TestingOverride()
}

func b(t *testing.T, address string) []byte {
	out, err := hex.DecodeString(address)
	require.NoError(t, err)

	return out
}

func bigString(t *testing.T, in string) *big.Int {
	v, ok := new(big.Int).SetString(in, 10)
	if !ok {
		t.Errorf("unable to convert string %q to big int", in)
	}
	return v
}
