package signer

import (
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/btcec"
	"github.com/streamingfast/eth-go"
)

// RecoverSigner recovers the signer of the received `signature` assuming that
// message `hash` was the message signed. If your message was signed using `SignPersonalHash`,
// use `RecoverPersonalSigner` instead which first re-compute the personal signing message.
//
// The `signature` must be in compact format as output by `SignHash` or
// `SignPersonalHash` which is the compact form `(r, s, v)` where `r` a 32
// bytes point, `s` is a second 32 bytes and `v` is the parity bit that will
// be either `27` or `28` where in hexadecimal compact form is the string
// `cfc0c9160c1fbe884f02298b76e194611904a4f18b6814dfd22f95b659b2f90b2564e455543316f125c3c51ec72b22917401d4872ddabb9a7a2d0bea1a827db31b`.
func RecoverSigner(signature eth.Hex, hash eth.Hash) (eth.Address, error) {
	btcecSignature := make([]byte, len(signature))

	btcecSignature[0] = signature[len(signature)-1]
	copy(btcecSignature[1:], signature[0:len(signature)-1])

	publicKey, compressed, err := btcec.RecoverCompact(btcec.S256(), btcecSignature, hash)
	if err != nil {
		return nil, fmt.Errorf("internal recover compact: %w", err)
	}

	// Original key was compressed, is it possible in our usage? For now, just ignore it
	_ = compressed

	return eth.NewPublicKeyFromECDSA(publicKey.ToECDSA()).Address(), nil
}

// RecoverPersonalSigner computes the signing message hash according to ERC-712 rules
// and then call `RecoverSigner` with the compute personal message hash.
func RecoverPersonalSigner(signature eth.Hex, hash eth.Hash) (eth.Address, error) {
	return RecoverSigner(signature, computePersonalMessageHash(hash))
}

var messagePrefix = []byte("\x19Ethereum Signed Message:\n")

func computePersonalMessageHash(hash eth.Hash) eth.Hash {
	lengthString := strconv.FormatUint(uint64(len(hash)), 10)
	data := make([]byte, len(messagePrefix)+len(lengthString)+len(hash))

	copy(data, messagePrefix)
	copy(data[len(messagePrefix):], []byte(lengthString))
	copy(data[len(messagePrefix)+len(lengthString):], hash)

	return eth.Keccak256(data)
}
