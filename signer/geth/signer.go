package geth

import (
	"fmt"
	"math/big"

	"github.com/dfuse-io/eth-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

type PrivateKeySigner struct {
	privateKey *eth.PrivateKey
	signer     types.Signer
	logger     *zap.Logger
}

func NewPrivateKeySigner(logger *zap.Logger, chainID uint64, privateKey *eth.PrivateKey) (*PrivateKeySigner, error) {
	return &PrivateKeySigner{
		privateKey: privateKey,
		signer:     types.NewEIP155Signer(big.NewInt(int64(chainID))),
		logger:     logger,
	}, nil
}

func (s *PrivateKeySigner) Sign(nonce uint64, to []byte, value *big.Int, gasLimit uint64, gasPrice *big.Int, trxData []byte) ([]byte, error) {
	s.logger.Info("signing transaction",
		zap.Uint64("nonce", nonce),
		zap.Stringer("to", eth.Address(to)),
		zap.Stringer("value", value),
		zap.Uint64("gas_limit", gasLimit),
		zap.Stringer("gas_price", gasPrice),
		zap.Stringer("trx_data", eth.Hex(trxData)),
	)

	tx := types.NewTransaction(nonce, common.BytesToAddress(to), value, gasLimit, gasPrice, trxData)
	signedTx, err := types.SignTx(tx, s.signer, s.privateKey.ToECDSA())
	if err != nil {
		return nil, fmt.Errorf("unable to sign transaction: %w", err)
	}

	ts := types.Transactions{signedTx}

	return ts.GetRlp(0), nil
}
