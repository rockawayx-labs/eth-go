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

package native

import (
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/streamingfast/eth-go"
	"github.com/streamingfast/eth-go/rlp"
	"go.uber.org/zap"
)

var b2 = big.NewInt(2)
var b35 = big.NewInt(35)
var b36 = big.NewInt(36)

// PrivateKeySigner respects EIP-155 signing rules regarding the exacy payload constructed
// from the Transation data.
type PrivateKeySigner struct {
	chainID        *big.Int
	chainIDDoubled *big.Int
	privateKey     *eth.PrivateKey
	logger         *zap.Logger
}

func NewPrivateKeySigner(logger *zap.Logger, chainID *big.Int, privateKey *eth.PrivateKey) (*PrivateKeySigner, error) {
	return &PrivateKeySigner{
		chainID:        chainID,
		chainIDDoubled: new(big.Int).Mul(chainID, b2),
		privateKey:     privateKey,
		logger:         logger,
	}, nil
}

func (p *PrivateKeySigner) Sign(nonce uint64, to []byte, value *big.Int, gasLimit uint64, gasPrice *big.Int, trxData []byte) (signedEncodedTrx []byte, err error) {
	v, r, s, err := p.Signature(nonce, to, value, gasLimit, gasPrice, trxData)
	if err != nil {
		return nil, err
	}

	p.logger.Debug("signed transaction signature",
		zap.Stringer("v", v),
		zap.Stringer("r", r),
		zap.Stringer("s", s),
	)

	// FIXME: Inefficent, the Signature process already had to encode nonce, gasPrice, gasLimit, to, value and trxData.
	//        We need to "pop" chainID, 0, 0 and replaced those by v, r, s values instead. Just trying to encode the generic
	//        then append the specific part does not work, prefixes change when doing and encoding does not work. So we need
	//        some kind of state. Maybe tweaking a bit the RLP encoding scheme could work here.
	data, err := rlp.Encode([]interface{}{
		nonce,
		gasPrice,
		gasLimit,
		to,
		value,
		trxData,
		v,
		r,
		s,
	})
	if err != nil {
		return nil, fmt.Errorf("rlp signed encode: %w", err)
	}

	return data, nil
}

func (p *PrivateKeySigner) Signature(nonce uint64, to []byte, value *big.Int, gasLimit uint64, gasPrice *big.Int, trxData []byte) (v, r, s *big.Int, err error) {
	p.logger.Debug("signing transaction",
		zap.Uint64("nonce", nonce),
		zap.Stringer("to", eth.Address(to)),
		zap.Stringer("value", value),
		zap.Uint64("gas_limit", gasLimit),
		zap.Stringer("gas_price", gasPrice),
		zap.Stringer("trx_data", eth.Hex(trxData)),
		zap.Stringer("chain_id", p.chainID),
	)

	data, err := rlp.Encode([]interface{}{
		nonce,
		gasPrice,
		gasLimit,
		to,
		value,
		trxData,
		p.chainID,
		uint64(0),
		uint64(0),
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("rlp encode: %w", err)
	}

	hash := eth.Keccak256(data)

	privKey := (*btcec.PrivateKey)(p.privateKey.ToECDSA())
	compressedSignature, err := btcec.SignCompact(btcec.S256(), privKey, hash, false)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("sign compact: %w", err)
	}

	r = new(big.Int).SetBytes(compressedSignature[1:33])
	s = new(big.Int).SetBytes(compressedSignature[33:])

	// In btcec, a `v` (i.e. byte at [0]) of 27 means the parity value of Y was 0
	// and if the parity was 1, the value will be 28. In Ethereum EIP-155, we are
	// looking for ChainID * 2 + 35 if parity is 0 and ChainID * 2 + 36 if parity
	// is 1.
	//
	// So, we first determined our "fixed" recovery ID either 35 or 36 based on the
	// btcec `v` value.
	recoveryFixedID := b35
	if compressedSignature[0] == 28 {
		recoveryFixedID = b36
	}

	// Then we apply the v = ChainID * 2 + {35, 36} math
	v = new(big.Int).Add(p.chainIDDoubled, recoveryFixedID)

	return v, r, s, nil
}
