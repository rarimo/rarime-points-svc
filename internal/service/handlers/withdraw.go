package handlers

import (
	"fmt"
	"net/http"
	"time"

	cosmos "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func Withdraw(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewWithdraw(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(req.Data.ID)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	balance, err := getBalanceByDID(req.Data.ID, true, r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if balance == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if err := isEligibleToWithdraw(balance, req.Data.Attributes.Amount); err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var withdrawal *data.Withdrawal
	err = EventsQ(r).Transaction(func() error {
		err = BalancesQ(r).FilterByDID(req.Data.ID).UpdateAmountBy(-req.Data.Attributes.Amount)
		if err != nil {
			return fmt.Errorf("decrease points amount: %w", err)
		}

		withdrawal, err = WithdrawalsQ(r).Insert(data.Withdrawal{
			UserDID: req.Data.ID,
			Amount:  req.Data.Attributes.Amount,
			Address: req.Data.Attributes.Address,
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
		Log(r).WithError(err).Error("Failed to perform withdrawal")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// balance should exist cause of previous logic
	balance, err = getBalanceByDID(req.Data.ID, true, r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID")
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
				ID:   balance.DID,
				Type: resources.BALANCE,
			},
		},
	}

	resp := resources.WithdrawalResponse{Data: wm}
	bm := newBalanceModel(balance)
	resp.Included.Add(&bm)

	return &resp
}

func isEligibleToWithdraw(balance *data.Balance, amount int64) validation.Errors {
	mapValidationErr := func(field, format string, a ...any) validation.Errors {
		return validation.Errors{
			field: fmt.Errorf(format, a...),
		}
	}

	if !balance.PassportHash.Valid {
		return mapValidationErr("is_verified", "user must have verified passport for withdrawals")
	}
	if balance.PassportExpires.Time.Before(time.Now().UTC()) {
		return mapValidationErr("is_verified", "user passport is expired")
	}
	if balance.Amount < amount {
		return mapValidationErr("data/attributes/amount", "insufficient balance: %d", balance.Amount)
	}

	return nil
}

func broadcastWithdrawalTx(req resources.WithdrawRequest, r *http.Request) error {
	urmo := req.Data.Attributes.Amount * PointPrice(r)
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
