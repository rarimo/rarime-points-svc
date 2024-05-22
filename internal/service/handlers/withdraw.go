package handlers

import (
	"fmt"
	"net/http"
	"time"

	cosmos "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/decentralized-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func Withdraw(w http.ResponseWriter, r *http.Request) {

	if PointPrice(r).Disabled {
		w.WriteHeader(http.StatusNoContent)
		return
	}

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

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(req.Data.ID)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	balance, err := BalancesQ(r).FilterByNullifier(req.Data.ID).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get balance by nullifier")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if balance == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if err = isEligibleToWithdraw(balance, req.Data.Attributes.Amount); err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var withdrawal *data.Withdrawal
	err = EventsQ(r).Transaction(func() error {
		err = BalancesQ(r).FilterByNullifier(req.Data.ID).UpdateAmountBy(-req.Data.Attributes.Amount)
		if err != nil {
			return fmt.Errorf("decrease points amount: %w", err)
		}

		withdrawal, err = WithdrawalsQ(r).Insert(data.Withdrawal{
			Nullifier: req.Data.ID,
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
	balance, err = BalancesQ(r).GetWithRank(req.Data.ID)
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

func isEligibleToWithdraw(balance *data.Balance, amount int64) error {
	mapValidationErr := func(field, format string, a ...any) validation.Errors {
		return validation.Errors{
			field: fmt.Errorf(format, a...),
		}
	}

	switch {
	case !balance.ReferredBy.Valid:
		return mapValidationErr("is_disabled", "user must be referred to withdraw")
	case !balance.PassportHash.Valid:
		return mapValidationErr("is_verified", "user must have verified passport to withdraw")
	case balance.PassportExpires.Time.Before(time.Now().UTC()):
		return mapValidationErr("is_verified", "user passport is expired")
	case !balance.IsWithdrawalAllowed:
		return mapValidationErr("is_withdrawal_allowed", "withdrawal ability was disabled for this user")
	case balance.Amount < amount:
		return mapValidationErr("data/attributes/amount", "insufficient balance: %d", balance.Amount)
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
