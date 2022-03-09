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

package signer

import (
	"math/big"

	"github.com/streamingfast/eth-go"
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

	// SignHash generates the signature for the according message hash. The signature returns is the compact version in the form
	// `(r, s, v)` where `r` a 32 bytes point, `s` is a second 32 bytes and `v` is the parity bit that will be either `27` or `28`.
	SignHash(hash eth.Hash) (signature []byte, err error)

	// SignPersonalHash generates the signature for the according message `hash` using the [ERC-712](https://eips.ethereum.org/EIPS/eip-712)
	// rules which is to hash `specific string`, `length` of message and `message` bytes than sign that element.
	SignPersonalHash(hash eth.Hash) (signature []byte, err error)

	// Signature generates the right payload for signing, perform the signing operation, extract the signature (v, r, s) and
	// return them.
	Signature(nonce uint64, toAddress []byte, value *big.Int, gasLimit uint64, gasPrice *big.Int, transactionData []byte) (r, s, v *big.Int, err error)
}
