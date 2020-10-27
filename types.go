package eth

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

type Hex []byte

func (h Hex) String() string {
	return hex.EncodeToString([]byte(h))
}

type Hash []byte

func (h Hash) String() string {
	return hex.EncodeToString([]byte(h))
}

type Address []byte

func MustNewAddress(input string) Address {
	input = strings.TrimPrefix(input, "0x")
	out, err := NewAddress(input)
	if err != nil {
		panic(fmt.Errorf("unable to create address: %w", err))
	}
	return out
}

func NewAddress(input string) (out Address, err error) {
	input = SanitizeHex(input)
	bytes, err := hex.DecodeString(input)
	if err != nil {
		return out, fmt.Errorf("invalid address: %w", err)
	}

	byteCount := len(bytes)
	if byteCount > 20 {
		bytes = bytes[byteCount-20:]
	}

	return bytes, nil
}

func (a Address) String() string {
	return hex.EncodeToString(a)
}

func (a Address) Pretty() string {
	return "0x" + hex.EncodeToString(a)
}

func (a Address) Bytes() []byte {
	return a[:]
}

func (a Address) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(a[:])), nil
}

func (a Address) ID() uint64 {
	return binary.LittleEndian.Uint64(a)
}

func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString([]byte(a)))
}

func (a *Address) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	var err error
	if *a, err = hex.DecodeString(strings.TrimPrefix(s, "0x")); err != nil {
		return err
	}

	return nil
}
