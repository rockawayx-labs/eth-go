package eth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestABIContract_Parse(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		expected    *ABI
		expectedErr error
	}{
		{
			"log event indexed",
			`[{"name":"PairCreated","type":"event","inputs":[{"indexed":true,"internalType":"address","name":"token0","type":"address"}]}]`,
			&ABI{
				MethodsMap: map[string]*MethodDef{},
				LogEventsMap: map[string]*LogEventDef{
					string(b(t, "b14a725aeeb25d591b81b16b4c5b25403dd8867bdd1876fa787867f566206be1")): {
						Name: "PairCreated",
						Parameters: []*LogParameter{
							{Name: "token0", TypeName: "address", Indexed: true},
						},
					},
				},
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			abi, err := parseABIFromBytes([]byte(test.in))
			if test.expectedErr == nil {
				require.NoError(t, err)
				abiEquals(t, test.expected, abi)
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func TestABIContract_ParseFile(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		expectedErr error
	}{
		{"standard", "testdata/uniswap_v2_factory.abi.json", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseABI(test.filename)
			if test.expectedErr == nil {
				require.NoError(t, err)
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func abiEquals(t *testing.T, expected *ABI, actual *ABI) {
	if len(expected.LogEventsMap) != len(actual.LogEventsMap) {
		require.Equal(t, expected.LogEventsMap, actual.LogEventsMap)
	} else {
		for key, value := range expected.LogEventsMap {
			assert.Contains(t, actual.LogEventsMap, key, "log event id %x", []byte(key))
			assert.Equal(t, value, actual.LogEventsMap[key])
		}
	}

	if len(expected.MethodsMap) != len(actual.MethodsMap) {
		require.Equal(t, expected.MethodsMap, actual.MethodsMap)
	} else {
		for key, value := range expected.MethodsMap {
			assert.Contains(t, actual.MethodsMap, key, "method id %x", []byte(key))
			assert.Equal(t, value, actual.MethodsMap[key])
		}
	}
}
