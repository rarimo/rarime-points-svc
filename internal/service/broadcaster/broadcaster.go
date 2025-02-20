package broadcaster

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rarimo/rarime-points-svc/internal/config"
	"gitlab.com/distributed_lab/logan/v3"
)

type Broadcaster struct {
	log         *logan.Entry
	broadcaster config.Broadcaster
}

func New(broadcaster config.Broadcaster, log *logan.Entry) (*Broadcaster, error) {
	return &Broadcaster{
		log:         log,
		broadcaster: broadcaster,
	}, nil
}

func (b *Broadcaster) BroadcastTransfer(ctx context.Context, to common.Address, amount *big.Int) error {
	b.broadcaster.LockNonce()
	defer b.broadcaster.UnlockNonce()

	nonce, err := b.broadcaster.GetCurrentNonce(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current nonce: %w", err)
	}

	gasPrice, err := b.broadcaster.RPC.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}
	gasPriceMax := b.broadcaster.MultiplyGasPrice(gasPrice)

	// Prepare transfer data
	transferFnSignature := []byte("transfer(address,uint256)")
	methodID := crypto.Keccak256(transferFnSignature)[:4]

	paddedAddress := common.LeftPadBytes(to.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// Estimate gas
	gasLimit, err := b.broadcaster.RPC.EstimateGas(ctx, ethereum.CallMsg{
		From: b.broadcaster.PublicKey,
		To:   &b.broadcaster.ERC20Contract,
		Data: data,
	})
	if err != nil {
		return fmt.Errorf("failed to estimate gas: %w", err)
	}

	// Create transaction
	tx := &types.DynamicFeeTx{
		ChainID:   b.broadcaster.ChainID,
		Nonce:     nonce,
		GasTipCap: gasPrice,    // priority fee
		GasFeeCap: gasPriceMax, // max fee set to multiplied gas price
		Gas:       gasLimit,
		To:        &b.broadcaster.ERC20Contract,
		Value:     big.NewInt(0),
		Data:      data,
	}

	signedTx, err := types.SignNewTx(b.broadcaster.PrivateKey, types.LatestSignerForChainID(b.broadcaster.ChainID), tx)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	err = b.broadcaster.RPC.SendTransaction(ctx, signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	b.broadcaster.IncrementNonce()

	b.log.WithFields(logan.F{
		"tx_hash": signedTx.Hash().Hex(),
		"to":      to.Hex(),
		"amount":  amount.String(),
		"nonce":   nonce,
	}).Info("ERC20 transfer transaction sent")

	return nil
}
