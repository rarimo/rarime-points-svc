package config

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	facesmt "github.com/rarimo/rarime-points-svc/tests/mocked/face_smt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	zkptypes "github.com/iden3/go-rapidsnark/types"
	zkpverifier "github.com/iden3/go-rapidsnark/verifier"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

const faceVerificationKey = "./verification_key.json"

var (
	ErrInvalidRoot = errors.New("invalid root")
)

const (
	FaceNullifierTreeRoot = iota
	FaceChallengedNullifier
)

type FaceVerifierer interface {
	FaceVerifier() *FaceVerifier
}

func NewFaceVerifier(getter kv.Getter) FaceVerifierer {
	return &faceVerifier{
		getter: getter,
	}
}

type faceVerifier struct {
	once   comfig.Once
	getter kv.Getter
}

type FaceVerifier struct {
	RPC              *ethclient.Client `fig:"rpc,required"`
	FaceStateAddress common.Address    `fig:"contract,required"`

	verificationKey []byte
}

func (c *faceVerifier) FaceVerifier() *FaceVerifier {
	return c.once.Do(func() interface{} {

		var cfg FaceVerifier

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "face_verifier")).
			With(figure.EthereumHooks, figure.BaseHooks).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out face verifier config: %w", err))
		}

		cfg.verificationKey, err = os.ReadFile(faceVerificationKey)
		if err != nil {
			panic(fmt.Errorf("failed to read faceVerificationKey: %w", err))
		}

		return &cfg
	}).(*FaceVerifier)
}

func (v *FaceVerifier) VerifyProof(proof zkptypes.ZKProof) error {
	root := decimalTo32Bytes(proof.PubSignals[FaceNullifierTreeRoot])
	if root == [32]byte{} {
		return ErrInvalidRoot
	}

	faceSMTCaller, err := facesmt.NewFaceSMTFiltererMock(v.FaceStateAddress, v.RPC)
	if err != nil {
		return fmt.Errorf("failed to create face smt caller: %w", err)
	}

	latestBlock, err := v.RPC.BlockNumber(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}

	it, err := faceSMTCaller.FilterRootUpdated(&bind.FilterOpts{
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
