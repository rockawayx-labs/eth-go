package eth

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

// KeyBag holds private keys in memory, for signing transactions.
type KeyBag struct {
	Keys []*PrivateKey `json:"keys"`
}

func NewKeyBag() *KeyBag {
	return &KeyBag{
		Keys: make([]*PrivateKey, 0),
	}
}

type PublicKey struct {
	inner ecdsa.PublicKey
}

func (p PublicKey) Address() Address {
	return pubkeyToAddress(p.inner)
}

type PrivateKey struct {
	inner *ecdsa.PrivateKey
}

func (p *PrivateKey) String() string {
	d := crypto.FromECDSA(p.inner)
	return hex.EncodeToString(d)
}

func (p *PrivateKey) Bytes() []byte {
	return p.inner.D.Bytes()
}

func (p *PrivateKey) ToECDSA() *ecdsa.PrivateKey {
	return p.inner
}

func (p *PrivateKey) MarshalJSON() ([]byte, error) {
	d := crypto.FromECDSA(p.inner)
	return json.Marshal(hex.EncodeToString(d))
}

func (p *PrivateKey) UnmarshalJSON(v []byte) (err error) {
	var s string
	if err = json.Unmarshal(v, &s); err != nil {
		return
	}

	newPrivKey, err := NewPrivateKey(s)
	if err != nil {
		return
	}

	*p = *newPrivKey

	return
}

func NewRandomPrivateKey() (*PrivateKey, error) {
	pkey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("unable to generate private key: %w", err)
	}
	return &PrivateKey{inner: pkey}, nil

}

func NewPrivateKey(rawPrivateKey string) (*PrivateKey, error) {
	pkey, err := crypto.HexToECDSA(rawPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to decode private key: %w", err)
	}
	return &PrivateKey{inner: pkey}, nil

}

func (p *PrivateKey) PublicKey() *PublicKey {
	return &PublicKey{inner: p.inner.PublicKey}
}

type keccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := sha3.NewLegacyKeccak256().(keccakState)
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}

func pubkeyToAddress(p ecdsa.PublicKey) Address {
	pubBytes := crypto.FromECDSAPub(&p)
	return Address(Keccak256(pubBytes[1:])[12:])
}
