package config

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	rootsmt "github.com/rarimo/rarime-points-svc/tests/mocked/root_smt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	zkptypes "github.com/iden3/go-rapidsnark/types"
	zkpverifier "github.com/iden3/go-rapidsnark/verifier"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

var (
	ErrInvalidRoot = errors.New("invalid root")
)

const (
	PubSignalNullifierTreeRoot = iota
	PubSignalChallengedNullifier
)

type RootInclusionVerifierer interface {
	RootInclusionVerifier() *RootInclusionVerifier
}

func NewRootInclusionVerifier(getter kv.Getter) RootInclusionVerifierer {
	return &rootVerifier{
		getter: getter,
	}
}

type rootVerifier struct {
	once   comfig.Once
	getter kv.Getter
}

type RootInclusionVerifier struct {
	RPC                 *ethclient.Client `fig:"rpc,required"`
	RootSMTAddress      common.Address    `fig:"contract,required"`
	VerificationKeyPath string            `fig:"verification_key_path,required"`

	verificationKey []byte
}

func (c *rootVerifier) RootInclusionVerifier() *RootInclusionVerifier {
	return c.once.Do(func() interface{} {

		var cfg RootInclusionVerifier

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "root_inclusion_verifier")).
			With(figure.EthereumHooks, figure.BaseHooks).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out root inclusion verifier config: %w", err))
		}

		cfg.verificationKey, err = os.ReadFile(cfg.VerificationKeyPath)
		if err != nil {
			panic(fmt.Errorf("failed to read verification key: %w", err))
		}

		return &cfg
	}).(*RootInclusionVerifier)
}

func (v *RootInclusionVerifier) VerifyProof(proof zkptypes.ZKProof) error {
	root := decimalTo32Bytes(proof.PubSignals[PubSignalNullifierTreeRoot])
	if root == [32]byte{} {
		return ErrInvalidRoot
	}

	rootSMTCaller, err := rootsmt.NewRootSMTFiltererMock(v.RootSMTAddress, v.RPC)
	if err != nil {
		return fmt.Errorf("failed to create root inclusion smt caller: %w", err)
	}

	latestBlock, err := v.RPC.BlockNumber(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}

	it, err := rootSMTCaller.FilterRootUpdated(&bind.FilterOpts{
		Start: max(0, latestBlock-5000),
	}, [][32]byte{root})
	if err != nil {
		return fmt.Errorf("failed to get root: %w", err)
	}

	if ok := it.Next(); !ok {
		return ErrInvalidRoot
	}

	if err = zkpverifier.VerifyGroth16(proof, v.verificationKey); err != nil {
		return fmt.Errorf("failed to verify proof: %w", err)
	}

	return nil
}

func decimalTo32Bytes(root string) [32]byte {
	b, ok := new(big.Int).SetString(root, 10)
	if !ok {
		return [32]byte{}
	}

	var bytes [32]byte
	b.FillBytes(bytes[:])

	return bytes
}
