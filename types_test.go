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
	"errors"
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddress_New(t *testing.T) {
	testNew(t, func(in string) (fmt.Stringer, error) { return NewAddress(in) })
}

func TestHex_New(t *testing.T) {
	testNew(t, func(in string) (fmt.Stringer, error) { return NewHex(in) })
}

func TestHash_New(t *testing.T) {
	testNew(t, func(in string) (fmt.Stringer, error) { return NewHash(in) })
}

func testNew(t *testing.T, new func(in string) (fmt.Stringer, error)) {
	tests := []struct {
		name        string
		in          string
		expected    string
		expectedErr error
	}{
		{"standard", "0xab", "ab", nil},
		{"odd length", "0xa", "0a", nil},
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
		{"odd length", `"770"`, "0770", nil},
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

func TestUint8_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		inJSON      string
		expected    uint8
		expectedErr error
	}{
		{"empty", `""`, 0, nil},
		{"hex empty", `"0x"`, 0, nil},
		{"hex odd length", `"0x1"`, 1, nil},
		{"hex mixed case", `"0xAc"`, 172, nil},
		{"hex boundary", `"0xff"`, math.MaxUint8, nil},
		{"hex outside boundary", `"0x0100"`, 0, errors.New(`invalid hex uint8 number: strconv.ParseUint: parsing "0100": value out of range`)},
		{"decimal", `"101"`, 101, nil},
		{"decimal boundary", `"255"`, math.MaxUint8, nil},
		{"decimal outside boundary", `"256"`, 0, errors.New(`invalid uint8 number: strconv.ParseUint: parsing "256": value out of range`)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var value Uint8
			err := json.Unmarshal([]byte(test.inJSON), &value)
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expected, uint8(value))
			} else {
				assert.EqualError(t, err, test.expectedErr.Error())
			}
		})
	}
}

func TestUint32_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		inJSON      string
		expected    uint32
		expectedErr error
	}{
		{"empty", `""`, 0, nil},
		{"hex empty", `"0x"`, 0, nil},
		{"hex odd length", `"0x1"`, 1, nil},
		{"hex mixed case", `"0x0AbC"`, 2748, nil},
		{"hex boundary", `"0xffffffff"`, math.MaxUint32, nil},
		{"hex outside boundary", `"0x0100000000"`, 0, errors.New(`invalid hex uint32 number: strconv.ParseUint: parsing "0100000000": value out of range`)},
		{"decimal", `"101"`, 101, nil},
		{"decimal boundary", `"4294967295"`, math.MaxUint32, nil},
		{"decimal outside boundary", `"4294967296"`, 0, errors.New(`invalid uint32 number: strconv.ParseUint: parsing "4294967296": value out of range`)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var value Uint32
			err := json.Unmarshal([]byte(test.inJSON), &value)
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expected, uint32(value))
			} else {
				assert.EqualError(t, err, test.expectedErr.Error())
			}
		})
	}
}

func TestUint64_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		inJSON      string
		expected    uint64
		expectedErr error
	}{
		{"empty", `""`, 0, nil},
		{"hex empty", `"0x"`, 0, nil},
		{"hex odd length", `"0x1"`, 1, nil},
		{"hex mixed case", `"0x0AbC"`, 2748, nil},
		{"hex boundary", `"0xffffffffffffffff"`, math.MaxUint64, nil},
		{"hex outside boundary", `"0x010000000000000000"`, 0, errors.New(`invalid hex uint64 number: strconv.ParseUint: parsing "010000000000000000": value out of range`)},
		{"decimal", `"101"`, 101, nil},
		{"decimal boundary", `"18446744073709551615"`, math.MaxUint64, nil},
		{"decimal outside boundary", `"18446744073709551616"`, 0, errors.New(`invalid uint64 number: strconv.ParseUint: parsing "18446744073709551616": value out of range`)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var value Uint64
			err := json.Unmarshal([]byte(test.inJSON), &value)
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expected, uint64(value))
			} else {
				assert.EqualError(t, err, test.expectedErr.Error())
			}
		})
	}
}

func Test_padTo32Bytes(t *testing.T) {
	type args struct {
		in []byte
	}

	tests := []struct {
		name    string
		args    args
		wantOut *Topic
	}{
		{
			"empty",
			args{in: nil},
			topic("0x0000000000000000000000000000000000000000000000000000000000000000"),
		},
		{
			"address",
			args{in: MustNewAddress("0xFffDB7377345371817F2b4dD490319755F5899eC")},
			topic("0x000000000000000000000000FffDB7377345371817F2b4dD490319755F5899eC"),
		},
		{
			"flush",
			args{in: MustNewHash("0x1111111111111111111111111111111111111111111111111111111111111111")},
			topic("0x1111111111111111111111111111111111111111111111111111111111111111"),
		},
		{
			"over",
			args{in: MustNewHash("0x111111111111111111111111111111111111111111111111111111111111111100000000")},
			topic("0x1111111111111111111111111111111111111111111111111111111111111111"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := padToTopic(tt.args.in); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("padTo32Bytes() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func topic(in string) (out *Topic) {
	var bytes [32]byte
	copy(bytes[:], MustNewHash(in))
	return (*Topic)(&bytes)
}
