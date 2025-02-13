package config

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/dig"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

const broadcasterYamlKey = "broadcaster"

type Broadcaster struct {
	RPC           *ethclient.Client
	ChainID       *big.Int
	PrivateKey    *ecdsa.PrivateKey
	PublicKey     common.Address
	ERC20Contract common.Address

	gasMultiplier float64
	nonce         uint64
	mut           *sync.Mutex
	lastUpdate    time.Time
}

type Broadcasterer interface {
	Broadcaster() Broadcaster
}

type broadcasterer struct {
	getter kv.Getter
	once   comfig.Once
}

func NewBroadcaster(getter kv.Getter) Broadcasterer {
	return &broadcasterer{
		getter: getter,
	}
}

func (b *broadcasterer) Broadcaster() Broadcaster {
	return b.once.Do(func() interface{} {
		var cfg struct {
			RPC              *ethclient.Client `fig:"rpc,required"`
			ChainID          *big.Int          `fig:"chain_id,required"`
			SenderPrivateKey *ecdsa.PrivateKey `fig:"sender_private_key"`
			ERC20Contract    common.Address    `fig:"erc20_contract,required"`
			GasMultiplier    float64           `fig:"gas_multiplier"`
		}

		err := figure.
			Out(&cfg).
			With(figure.BaseHooks, figure.EthereumHooks).
			From(kv.MustGetStringMap(b.getter, broadcasterYamlKey)).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out broadcaster: %w", err))
		}

		if cfg.SenderPrivateKey == nil {
			cfg.SenderPrivateKey = extractPubKey()
		}

		gasMultiplier := float64(1)
		if cfg.GasMultiplier > 0 {
			gasMultiplier = cfg.GasMultiplier
		}

		publicKey := crypto.PubkeyToAddress(cfg.SenderPrivateKey.PublicKey)
		nonce, err := cfg.RPC.NonceAt(context.Background(), publicKey, nil)
		if err != nil {
			panic(fmt.Errorf("failed to get nonce %w", err))
		}

		return Broadcaster{
			RPC:           cfg.RPC,
			PrivateKey:    cfg.SenderPrivateKey,
			PublicKey:     publicKey,
			ChainID:       cfg.ChainID,
			ERC20Contract: cfg.ERC20Contract,

			gasMultiplier: gasMultiplier,
			nonce:         nonce,
			mut:           &sync.Mutex{},
			lastUpdate:    time.Now().UTC(),
		}
	}).(Broadcaster)
}

func extractPubKey() *ecdsa.PrivateKey {
	var envPK struct {
		PrivateKey *ecdsa.PrivateKey `dig:"PRIVATE_KEY,clear"`
	}

	if err := dig.Out(&envPK).With(figure.EthereumHooks).Now(); err != nil {
		panic(fmt.Errorf("failed to figure out private key from ENV: %w", err))
	}

	return envPK.PrivateKey
}

func (n *Broadcaster) LockNonce() {
	n.mut.Lock()
}

func (n *Broadcaster) UnlockNonce() {
	n.mut.Unlock()
}

func (n *Broadcaster) Nonce() uint64 {
	return n.nonce
}

func (n *Broadcaster) IncrementNonce() {
	n.nonce++
}

func (n *Broadcaster) GetCurrentNonce(ctx context.Context) (uint64, error) {
	if time.Since(n.lastUpdate) > 30*time.Second {
		nonce, err := n.RPC.NonceAt(ctx, n.PublicKey, nil)
		if err != nil {
			return 0, fmt.Errorf("failed to get nonce: %w", err)
		}
		n.nonce = nonce
		n.lastUpdate = time.Now().UTC()
	}
	return n.nonce, nil
}

func (n *Broadcaster) MultiplyGasPrice(gasPrice *big.Int) *big.Int {
	var ONE = 1000000000 // ONE - One GWEI
	mult := big.NewFloat(0).Mul(big.NewFloat(n.gasMultiplier), big.NewFloat(float64(ONE)))
	gas, _ := big.NewFloat(0).Mul(big.NewFloat(0).SetInt(gasPrice), mult).Int(nil)
	return big.NewInt(0).Div(gas, big.NewInt(int64(ONE)))
}
