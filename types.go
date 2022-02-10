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
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Uint8 uint8

func (b *Uint8) UnmarshalText(text []byte) error {
	value, err := parseUint(string(text), 8)
	*b = Uint8(value)

	return err
}

type Uint16 uint16

func (b *Uint16) UnmarshalText(text []byte) error {
	value, err := parseUint(string(text), 16)
	*b = Uint16(value)

	return err
}

type Uint32 uint32

func (b *Uint32) UnmarshalText(text []byte) error {
	value, err := parseUint(string(text), 32)
	*b = Uint32(value)

	return err
}

type Uint64 uint64

func (b *Uint64) UnmarshalText(text []byte) error {
	value, err := parseUint(string(text), 64)
	*b = Uint64(value)

	return err
}

func parseUint(text string, bitSize int) (uint64, error) {
	if len(text) == 0 {
		return 0, nil
	}

	// If it's a hexadecimal string, let's parse it as-is
	if strings.HasPrefix(text, "0x") || strings.HasPrefix(text, "0X") {
		text = text[2:]
		if text == "" {
			return 0, nil
		}

		value, err := strconv.ParseUint(text, 16, bitSize)
		if err != nil {
			return 0, fmt.Errorf("invalid hex uint%d number: %w", bitSize, err)
		}

		return value, nil
	}

	// Otherwise, we assume it's a decimal number
	value, err := strconv.ParseUint(text, 10, bitSize)
	if err != nil {
		return 0, fmt.Errorf("invalid uint%d number: %w", bitSize, err)
	}

	return value, nil
}

type Hex []byte

func MustNewHex(input string) Hex {
	return Hex(mustNewByteSlice("hex", input))
}

func NewHex(input string) (Hex, error) {
	out, err := newByteSlice("hex", input)
	if err != nil {
		return nil, err
	}

	return Hex(out), nil
}

func (h Hex) String() string                   { return byteSlice(h).String() }
func (h Hex) Pretty() string                   { return byteSlice(h).Pretty() }
func (h Hex) Bytes() []byte                    { return h[:] }
func (h Hex) MarshalText() ([]byte, error)     { return byteSlice(h).MarshalText() }
func (h Hex) ID() uint64                       { return byteSlice(h).ID() }
func (h Hex) MarshalJSON() ([]byte, error)     { return byteSlice(h).MarshalJSON() }
func (h Hex) MarshalJSONRPC() ([]byte, error)  { return byteSlice(h).MarshalJSONRPC() }
func (h *Hex) UnmarshalJSON(data []byte) error { return (*byteSlice)(h).UnmarshalJSON(data) }

type Hash []byte

func MustNewHash(input string) Hash {
	return Hash(mustNewByteSlice("hash", input))
}

func NewHash(input string) (Hash, error) {
	out, err := newByteSlice("hash", input)
	if err != nil {
		return nil, err
	}

	return Hash(out), nil
}

func (h Hash) String() string                   { return byteSlice(h).String() }
func (h Hash) Pretty() string                   { return byteSlice(h).Pretty() }
func (h Hash) Bytes() []byte                    { return h[:] }
func (h Hash) MarshalText() ([]byte, error)     { return byteSlice(h).MarshalText() }
func (h Hash) ID() uint64                       { return byteSlice(h).ID() }
func (h Hash) MarshalJSON() ([]byte, error)     { return byteSlice(h).MarshalJSON() }
func (h Hash) MarshalJSONRPC() ([]byte, error)  { return byteSlice(h).MarshalJSONRPC() }
func (h *Hash) UnmarshalJSON(data []byte) error { return (*byteSlice)(h).UnmarshalJSON(data) }

type Address []byte

func MustNewAddress(input string) Address {
	out, err := NewAddress(input)
	if err != nil {
		panic(err)
	}

	return out
}

func NewAddress(input string) (Address, error) {
	out, err := newByteSlice("address", input)
	if err != nil {
		return nil, err
	}

	byteCount := len(out)
	if byteCount > 20 {
		out = out[byteCount-20:]
	}

	return Address(out), nil
}

func (a Address) String() string                   { return byteSlice(a).String() }
func (a Address) Pretty() string                   { return byteSlice(a).Pretty() }
func (a Address) Bytes() []byte                    { return a[:] }
func (a Address) MarshalText() ([]byte, error)     { return byteSlice(a).MarshalText() }
func (a Address) ID() uint64                       { return byteSlice(a).ID() }
func (a Address) MarshalJSON() ([]byte, error)     { return byteSlice(a).MarshalJSON() }
func (a Address) MarshalJSONRPC() ([]byte, error)  { return byteSlice(a).MarshalJSONRPC() }
func (a *Address) UnmarshalJSON(data []byte) error { return (*byteSlice)(a).UnmarshalJSON(data) }

type byteSlice []byte

func mustNewByteSlice(tag string, input string) byteSlice {
	out, err := newByteSlice(tag, input)
	if err != nil {
		panic(err)
	}

	return out
}

func newByteSlice(tag string, input string) (out byteSlice, err error) {
	bytes, err := hex.DecodeString(SanitizeHex(input))
	if err != nil {
		return out, fmt.Errorf("invalid %s %q: %w", tag, input, err)
	}

	return bytes, nil
}

func (b byteSlice) String() string {
	return hex.EncodeToString(b)
}

func (b byteSlice) Pretty() string {
	return "0x" + hex.EncodeToString(b)
}

func (b byteSlice) Bytes() []byte {
	return b
}

func (b byteSlice) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b byteSlice) ID() uint64 {
	return binary.LittleEndian.Uint64(b)
}

func (b byteSlice) MarshalJSON() ([]byte, error) {
	return []byte(`"` + hex.EncodeToString([]byte(b)) + `"`), nil
}

func (b byteSlice) MarshalJSONRPC() ([]byte, error) {
	return []byte(`"` + b.Pretty() + `"`), nil
}

func (b *byteSlice) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	s = strings.TrimPrefix(s, "0x")
	if len(s)%2 != 0 {
		s = "0" + s
	}

	var err error
	if *b, err = hex.DecodeString(s); err != nil {
		return err
	}

	return nil
}

type Topic [32]byte

func (f Topic) MarshalJSONRPC() ([]byte, error) {
	return []byte(`"0x` + hex.EncodeToString(f[:]) + `"`), nil
}

func LogTopic(in interface{}) *Topic {
	switch v := in.(type) {
	case string:
		return padToTopic(MustNewHash(v))
	case []byte:
		return padToTopic(v)
	case Hash:
		return padToTopic(v)
	case Hex:
		return padToTopic(v)
	case Address:
		return padToTopic(v)
	case nil:
		return nil
	default:
		valueOf := reflect.ValueOf(v)
		if valueOf.Kind() == reflect.Ptr && valueOf.IsNil() {
			return nil
		}

		panic(fmt.Errorf("don't know how to turn %T into a LogTopic", in))
	}
}

func padToTopic(in []byte) (out *Topic) {
	startOffset := 32 - len(in)
	if startOffset < 0 {
		startOffset = 0
	}

	var topic Topic
	copy(topic[startOffset:], in)

	return &topic
}
