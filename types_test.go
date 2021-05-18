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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddress_New(t *testing.T) {
	testNew(t, func(in string) (fmt.Stringer, error) { return NewAddress(in) })
}

func testNew(t *testing.T, new func(in string) (fmt.Stringer, error)) {
	tests := []struct {
		name        string
		in          string
		expected    string
		expectedErr error
	}{
		{"standard", "0xab", "ab", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := new(test.in)
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expected, value.String())
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func TestAddress_Pretty(t *testing.T) {
	testPretty(t, func(in []byte) string { return Address(in).Pretty() })
}

func TestHash_Pretty(t *testing.T) {
	testPretty(t, func(in []byte) string { return Hash(in).Pretty() })
}

func TestHex_Pretty(t *testing.T) {
	testPretty(t, func(in []byte) string { return Hex(in).Pretty() })
}

func testPretty(t *testing.T, pretty func(in []byte) string) {
	tests := []struct {
		name     string
		in       []byte
		expected string
	}{
		{"standard", []byte{0xab}, "0xab"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, pretty(test.in))
		})
	}
}

func TestAddress_UnmarshalJSON(t *testing.T) {
	testUnmarshalJSON(t, func(in []byte) (fmt.Stringer, error) {
		var out Hex
		return out, json.Unmarshal(in, &out)
	})
}

func TestHash_UnmarshalJSON(t *testing.T) {
	testUnmarshalJSON(t, func(in []byte) (fmt.Stringer, error) {
		var out Hash
		return out, json.Unmarshal(in, &out)
	})
}

func TestHex_UnmarshalJSON(t *testing.T) {
	testUnmarshalJSON(t, func(in []byte) (fmt.Stringer, error) {
		var out Hex
		return out, json.Unmarshal(in, &out)
	})
}

func testUnmarshalJSON(t *testing.T, unmarshalJSON func(jsonMessage []byte) (fmt.Stringer, error)) {
	tests := []struct {
		name        string
		inJSON      string
		expected    string
		expectedErr error
	}{
		{"standard", `"ab"`, "ab", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := unmarshalJSON([]byte(test.inJSON))
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expected, value.String())
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}
