package config

import (
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/rarimo/rarime-points-svc/tests/mocked/faceregistry"

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
)

const (
	PubSignalNullifier = iota
	PubSignalEventID
	PubSignalNonce
)

const likenessRegistryEventID = "00000000000000000000000000000000000000000000000000000000000000000000000000000"

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
	RPC                     *ethclient.Client `fig:"rpc,required"`
	RootSMTAddress          common.Address    `fig:"contract,required"`
	VerificationKeyPath     string            `fig:"verification_key_path,required"`
	LikenessRegistryEventID string            `fig:"likeness_registry_event_id,required"`

	verificationKey []byte
}

func (c *rootVerifier) RootInclusionVerifier() *RootInclusionVerifier {
	return c.once.Do(func() interface{} {

		cfg := RootInclusionVerifier{LikenessRegistryEventID: likenessRegistryEventID}

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
	nullifier, ok := new(big.Int).SetString(proof.PubSignals[PubSignalNullifier], 10)
	if !ok {
		return fmt.Errorf("failed to convert nullifier to *big.Int")
	}

	faceRegistryCaller, err := faceregistry.NewFaceRegistryCaller(v.RootSMTAddress, v.RPC)
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

	err = checkCmpBigIntFromStrings(proof.PubSignals[PubSignalEventID], v.LikenessRegistryEventID)
	if err != nil {
		return fmt.Errorf("failed to check event id: %w", err)
	}

	if err = zkpverifier.VerifyGroth16(proof, v.verificationKey); err != nil {
		return fmt.Errorf("failed to verify proof: %w", err)
	}

	return nil
}

func checkCmpBigIntFromStrings(a, b string) error {
	aInt, ok := new(big.Int).SetString(a, 10)
	if !ok {
		return fmt.Errorf("failed to convert %s to *big.Int", a)
	}
	bInt, ok := new(big.Int).SetString(b, 10)
	if !ok {
		return fmt.Errorf("failed to convert %s to *big.Int", b)
	}
	if aInt.Cmp(bInt) != 0 {
		return fmt.Errorf("is not equal")
	}
	return nil
}
