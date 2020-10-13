package eth

import (
	"encoding/binary"
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

type DecodedMethod struct {
	Signature  string
	Parameters map[string]interface{}
}

type LogEvent struct {
	Signature  string
	Parameters map[string]interface{}
}

type Log struct {
	Address []byte   `json:"address,omitempty"`
	Topics  [][]byte `json:"topics,omitempty"`
	Data    []byte   `json:"data,omitempty"`
	// supplement
	Index      uint32 `json:"index,omitempty"`
	BlockIndex uint32 `json:"blockIndex,omitempty"`
}
