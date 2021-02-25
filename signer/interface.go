package signer

import (
	"math/big"
)

// Signer is the interface of all implementation that are able to sign message data according to the various
// Ethereum rules.
//
// **Important** This interface might change at any time to adjust to new Ethereum rules.
type Signer interface {
	Sign(nonce uint64, toAddress []byte, value *big.Int, gasLimit uint64, gasPrice *big.Int, transactionData []byte) (r, s, v *big.Int, err error)
}
