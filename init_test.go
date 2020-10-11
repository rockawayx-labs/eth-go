package eth

import (
	"encoding/hex"
	"github.com/dfuse-io/logging"
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	logging.TestingOverride()
}

func b(t *testing.T, address string) []byte {
	out, err := hex.DecodeString(address)
	require.NoError(t, err)

	return out
}
