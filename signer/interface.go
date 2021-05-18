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
