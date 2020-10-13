package signer

import "math/big"

type Signer interface {
	Sign(nonce uint64, toAddress []byte, value *big.Int, gasLimit uint64, gasPrice *big.Int, transactionData []byte) ([]byte, error)
}
