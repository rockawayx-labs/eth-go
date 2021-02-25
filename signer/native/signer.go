package native

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/dfuse-io/eth-go"
	"github.com/dfuse-io/eth-go/rlp"
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

func NewPrivateKeySigner(logger *zap.Logger, chainID uint64, privateKey *eth.PrivateKey) (*PrivateKeySigner, error) {
	bigChainID := big.NewInt(int64(chainID))

	return &PrivateKeySigner{
		chainID:        bigChainID,
		chainIDDoubled: new(big.Int).Mul(bigChainID, b2),
		privateKey:     privateKey,
		logger:         logger,
	}, nil
}

func (p *PrivateKeySigner) Sign(nonce uint64, to []byte, value *big.Int, gasLimit uint64, gasPrice *big.Int, trxData []byte) (v, r, s *big.Int, err error) {
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

	fmt.Printf("Native (%s, %s, %s, %s)\n", v.String(), r.String(), s.String(), hex.EncodeToString(compressedSignature))
	return v, r, s, nil
}
