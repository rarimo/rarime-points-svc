package config

import (
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/rarimo/rarime-points-svc/internal/contracts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	zkptypes "github.com/iden3/go-rapidsnark/types"
	zkpverifier "github.com/iden3/go-rapidsnark/verifier"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

var (
	ErrUserNotRegistered = errors.New("user not registered in the face registry")
	maxEventID, _        = new(big.Int).SetString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
)

const (
	PubSignalNullifier = iota
	PubSignalEventID
	PubSignalNonce
)

type LikenessRegistryVerifierer interface {
	LikenessRegistryVerifier() *LikenessRegistryVerifier
}

func NewLikenessRegistryVerifier(getter kv.Getter) LikenessRegistryVerifierer {
	return &likenessVerifier{
		getter: getter,
	}
}

type likenessVerifier struct {
	once   comfig.Once
	getter kv.Getter
}

type LikenessRegistryVerifier struct {
	RPC                     *ethclient.Client `fig:"rpc,required"`
	LikenessContract        common.Address    `fig:"contract,required"`
	VerificationKeyPath     string            `fig:"verification_key_path,required"`
	LikenessRegistryEventID string            `fig:"likeness_registry_event_id,required"`

	verificationKey []byte
}

func (c *likenessVerifier) LikenessRegistryVerifier() *LikenessRegistryVerifier {
	return c.once.Do(func() interface{} {

		var cfg LikenessRegistryVerifier

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "likeness_registry_verifier")).
			With(figure.EthereumHooks, figure.BaseHooks).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out likeness registry verifier config: %w", err))
		}

		cfg.verificationKey, err = os.ReadFile(cfg.VerificationKeyPath)
		if err != nil {
			panic(fmt.Errorf("failed to read verification key: %w", err))
		}

		eventID, ok := new(big.Int).SetString(cfg.LikenessRegistryEventID, 10)
		if !ok {
			panic(fmt.Errorf("event_id must be valid decimal"))
		}

		if eventID.Cmp(maxEventID) == 1 {
			panic(fmt.Errorf("event_id must be less than 31 bytes"))
		}

		return &cfg
	}).(*LikenessRegistryVerifier)
}

func (v *LikenessRegistryVerifier) VerifyProof(proof zkptypes.ZKProof) error {
	nullifier, ok := new(big.Int).SetString(proof.PubSignals[PubSignalNullifier], 10)
	if !ok {
		return fmt.Errorf("failed to convert nullifier to *big.Int")
	}

	faceRegistryCaller, err := contracts.NewFaceRegistryCaller(v.LikenessContract, v.RPC)
	if err != nil {
		return fmt.Errorf("failed to create face registry caller: %w", err)
	}

	ok, err = faceRegistryCaller.IsUserRegistered(nil, nullifier)
	if err != nil {
		return fmt.Errorf("failed to check is user registered: %w", err)
	}
	if !ok {
		return ErrUserNotRegistered
	}

	if proof.PubSignals[PubSignalEventID] != v.LikenessRegistryEventID {
		return fmt.Errorf("invalid likeness registry event id")
	}

	if err = zkpverifier.VerifyGroth16(proof, v.verificationKey); err != nil {
		return fmt.Errorf("failed to verify proof: %w", err)
	}

	return nil
}
