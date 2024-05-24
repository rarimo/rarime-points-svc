package handlers

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	cosmos "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	zkptypes "github.com/iden3/go-rapidsnark/types"
	"github.com/rarimo/decentralized-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	zk "github.com/rarimo/zkverifier-kit"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const usaAuthorithy = "8571562"

func Withdraw(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewWithdraw(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	log := Log(r).WithFields(map[string]any{
		"nullifier":     req.Data.ID,
		"points_amount": req.Data.Attributes.Amount,
		"dest_address":  req.Data.Attributes.Address,
	})

	if PointPrice(r).Disabled {
		log.Debug("Withdrawal disabled!")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	nullifier := req.Data.ID

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(nullifier)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	balance, err := BalancesQ(r).FilterByNullifier(nullifier).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get balance by nullifier")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if balance == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var proof zkptypes.ZKProof
	if err := json.Unmarshal(req.Data.Attributes.Proof, &proof); err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// MustDecode will never panic, because of the previous logic
	proof.PubSignals[zk.Nullifier] = new(big.Int).SetBytes(hexutil.MustDecode(nullifier)).String()
	if err := Verifier(r).VerifyProof(proof, zk.WithProofSelectorValue("23073")); err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if proof.PubSignals[zk.Citizenship] == "usaAuthorithy" {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"authority": errors.New("Incorrect authority")})...)
		return
	}

	if err = isEligibleToWithdraw(r, balance, req.Data.Attributes.Amount); err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var withdrawal *data.Withdrawal
	err = EventsQ(r).Transaction(func() error {
		err = BalancesQ(r).FilterByNullifier(nullifier).UpdateAmountBy(-req.Data.Attributes.Amount)
		if err != nil {
			return fmt.Errorf("decrease points amount: %w", err)
		}

		withdrawal, err = WithdrawalsQ(r).Insert(data.Withdrawal{
			Nullifier: nullifier,
			Amount:    req.Data.Attributes.Amount,
			Address:   req.Data.Attributes.Address,
		})
		if err != nil {
			return fmt.Errorf("add withdrawal entry: %w", err)
		}

		if err = broadcastWithdrawalTx(req, r); err != nil {
			return fmt.Errorf("broadcast transfer tx: %w", err)
		}
		return nil
	})

	if err != nil {
		log.WithError(err).Error("Failed to perform withdrawal")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// balance should exist cause of previous logic
	balance, err = BalancesQ(r).GetWithRank(nullifier)
	if err != nil {
		log.WithError(err).Error("Failed to get balance by nullifier with rank")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newWithdrawResponse(*withdrawal, *balance))
}

func newWithdrawResponse(w data.Withdrawal, balance data.Balance) *resources.WithdrawalResponse {
	wm := newWithdrawalModel(w)
	wm.Relationships = &resources.WithdrawalRelationships{
		Balance: resources.Relation{
			Data: &resources.Key{
				ID:   balance.Nullifier,
				Type: resources.BALANCE,
			},
		},
	}

	resp := resources.WithdrawalResponse{Data: wm}
	bm := newBalanceModel(balance)
	resp.Included.Add(&bm)

	return &resp
}

func isEligibleToWithdraw(r *http.Request, balance *data.Balance, amount int64) error {
	mapValidationErr := func(field, format string, a ...any) validation.Errors {
		return validation.Errors{
			field: fmt.Errorf(format, a...),
		}
	}

	switch {
	case !balance.ReferredBy.Valid:
		return mapValidationErr("is_disabled", "user must be referred to withdraw")
	case balance.Amount < amount:
		return mapValidationErr("data/attributes/amount", "insufficient balance: %d", balance.Amount)
	case !Levels(r)[balance.Level].WithdrawalAllowed:
		return mapValidationErr("withdrawal not allowed", "user must up level to have withdraw ability")
	}

	return nil
}

func broadcastWithdrawalTx(req resources.WithdrawRequest, r *http.Request) error {
	urmo := req.Data.Attributes.Amount * PointPrice(r).PointPriceURMO
	tx := &bank.MsgSend{
		FromAddress: Broadcaster(r).Sender(),
		ToAddress:   req.Data.Attributes.Address,
		Amount:      cosmos.NewCoins(cosmos.NewInt64Coin("urmo", urmo)),
	}

	err := Broadcaster(r).BroadcastTx(r.Context(), tx)
	if err != nil {
		return fmt.Errorf("broadcast withdrawal tx: %w", err)
	}

	return nil
}
