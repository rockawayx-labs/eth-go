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

package rpc

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/streamingfast/eth-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalJSONRPC(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expected    string
		expectedErr error
	}{
		{"int8 1", int8(1), `"0x1"`, nil},
		{"int16 1", int16(1), `"0x1"`, nil},
		{"int32 1", int32(1), `"0x1"`, nil},
		{"int64 1", int64(1), `"0x1"`, nil},
		{"int 1", int(1), `"0x1"`, nil},

		{"uint8 1", uint8(1), `"0x1"`, nil},
		{"uint16 1", uint16(1), `"0x1"`, nil},
		{"uint32 1", uint32(1), `"0x1"`, nil},
		{"uint64 1", uint64(1), `"0x1"`, nil},
		{"uint 1", uint(1), `"0x1"`, nil},

		{"big.Int", *big.NewInt(1), `"0x1"`, nil},
		{"*big.Int", big.NewInt(1), `"0x1"`, nil},

		// Some wrapping tests
		{"wrapped no marshal RPC", wrappedUint64NoMarshalRPC(1), `"0x1"`, nil},
		{"wrapped no marshal RPC, custom UnmarshalJSON", wrappedUint64CustomUnmarshal(1), `"0x1"`, nil},

		// Quantity vs Data zero handling (quantities should be 0x0, data should be 0x)
		{"quantity int8 0", int8(0), `"0x0"`, nil},
		{"quantity int16 0", int16(0), `"0x0"`, nil},
		{"quantity int32 0", int32(0), `"0x0"`, nil},
		{"quantity int64 0", int64(0), `"0x0"`, nil},
		{"quantity int 0", int(0), `"0x0"`, nil},

		{"quantity uint8 0", uint8(0), `"0x0"`, nil},
		{"quantity uint16 0", uint16(0), `"0x0"`, nil},
		{"quantity uint32 0", uint32(0), `"0x0"`, nil},
		{"quantity uint64 0", uint64(0), `"0x0"`, nil},
		{"quantity uint 0", uint(0), `"0x0"`, nil},

		{"quantity bigInt nil", (*big.Int)(nil), `"0x0"`, nil},
		{"quantity bigInt 0", big.Int{}, `"0x0"`, nil},

		{"quantity trims zero prefix", int(15), `"0xf"`, nil},

		{"data []byte nil", ([]byte)(nil), `"0x"`, nil},
		{"data []byte empty", []byte{}, `"0x"`, nil},
		{"data []byte 0x00", []byte{0x00}, `"0x00"`, nil},

		{"*eth.MethodCall", eth.MustNewMethodDef("name()").NewCall(), `"0x06fdde03"`, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := MarshalJSONRPC(test.in)
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.JSONEq(t, test.expected, string(actual))
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

type wrappedUint64NoMarshalRPC uint64

type wrappedUint64CustomUnmarshal uint64

func (u *wrappedUint64CustomUnmarshal) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("invalid input: %s", string(data))
	}

	value, err := strconv.ParseUint(string(data[1:len(data)-1]), 10, 64)
	if err != nil {
		return err
	}

	*u = wrappedUint64CustomUnmarshal(value)
	return nil
}
