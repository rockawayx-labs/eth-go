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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestABIContract_Parse(t *testing.T) {
	tests := []struct {
		name        string
		expected    *ABI
		expectedErr error
	}{
		{
			"log event indexed",
			&ABI{
				FunctionsMap: map[string]*MethodDef{},
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

		{
			"struct tuple alone",
			&ABI{
				FunctionsMap: map[string]*MethodDef{
					string(b(t, "b961c32d")): {
						Name: "tupleAlone",
						Parameters: []*MethodParameter{
							{Name: "period", TypeName: "tuple", InternalType: "struct ClaimPeriod", TypeMutability: "", Components: []*StructComponent{
								{Name: "tokenID", Type: "uint256", InternalType: "uint256"},
								{Name: "fromBlockNum", Type: "uint64", InternalType: "uint64"},
								{Name: "toBlockNum", Type: "uint64", InternalType: "uint64"},
							}},
						},
						StateMutability: StateMutabilityPure,
					},
				},
			},
			nil,
		},

		{
			"struct tuple array",
			&ABI{
				FunctionsMap: map[string]*MethodDef{
					string(b(t, "15eb963d")): {
						Name: "tupleArray",
						Parameters: []*MethodParameter{
							{Name: "periods", TypeName: "tuple[]", InternalType: "struct ClaimPeriod[]", TypeMutability: "", Components: []*StructComponent{
								{Name: "tokenID", Type: "uint256", InternalType: "uint256"},
								{Name: "fromBlockNum", Type: "uint64", InternalType: "uint64"},
								{Name: "toBlockNum", Type: "uint64", InternalType: "uint64"},
							}},
						},
						StateMutability: StateMutabilityView,
					},
				},
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			abiFile := filepath.Join("testdata", strings.ReplaceAll(test.name, " ", "_")+".abi.json")
			stat, err := os.Stat(abiFile)
			require.NoError(t, err, "unable to stat abi file %q (inferred from %s space replaced by _)", abiFile, test.name)
			require.Equal(t, true, !stat.IsDir(), "abi file %q is a directory (inferred from %s space replaced by _)", abiFile, test.name)

			abi, err := ParseABI(abiFile)
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
			assert.Contains(t, actual.LogEventsMap, key, "log event id %s", key)
			assert.Equal(t, value, actual.LogEventsMap[key])
		}
	}

	if len(expected.FunctionsMap) != len(actual.FunctionsMap) {
		require.Equal(t, expected.FunctionsMap, actual.FunctionsMap)
	} else {
		for key, value := range expected.FunctionsMap {
			assert.Contains(t, actual.FunctionsMap, key, "method id %s", key)
			assert.Equal(t, value, actual.FunctionsMap[key])
		}
	}
}
