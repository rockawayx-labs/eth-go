package eth

import (
	"encoding/hex"
	"fmt"
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
	out, err := NewAddress(input)
	if err != nil {
		panic(fmt.Errorf("unable to create address: %w", err))
	}
	return out
}

func NewAddress(input string) (out Address, err error) {
	bytes, err := hex.DecodeString(input)
	if err != nil {
		return out, fmt.Errorf("invalid address: %w", err)
	}

	byteCount := len(bytes)
	if byteCount > 20 {
		bytes = bytes[byteCount-20:]
	}

	return Address(bytes), nil
}

func (a Address) String() string {
	return hex.EncodeToString(a)
}

func (a Address) Bytes() []byte {
	return a[:]
}

func (a Address) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(a[:])), nil
}

type DecodedMethod struct {
	Signature  string
	Parameters map[string]interface{}
}

type LogEvent struct {
	Signature  string
	Parameters map[string]interface{}
}
