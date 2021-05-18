package signer

import (
	"math/big"
)

// Signer is the interface of all implementation that are able to sign message data according to the various
// Ethereum rules.
//
// **Important** This interface might change at any time to adjust to new Ethereum rules.
type Signer interface {
	// Sign generates the right payload for signing, perform the signing operation, extract the signature (v, r, s) from it then
	// complete the standard Ethereum transaction signing process by appending r, s to transaction payload and completing the RLP
	// encoding.
	Sign(nonce uint64, toAddress []byte, value *big.Int, gasLimit uint64, gasPrice *big.Int, transactionData []byte) (signedEncodedTrx []byte, err error)

	// Signature generates the right payload for signing, perform the signing operation, extract the signature (v, r, s) and
	// return them.
	Signature(nonce uint64, toAddress []byte, value *big.Int, gasLimit uint64, gasPrice *big.Int, transactionData []byte) (r, s, v *big.Int, err error)
}
