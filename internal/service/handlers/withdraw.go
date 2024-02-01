package handlers

import (
	"fmt"
	"net/http"

	cosmos "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

	if !isEnoughPoints(req, w, r) {
		return
	}

	err = BalancesQ(r).FilterByDID(req.Data.ID).UpdateAmountBy(-req.Data.Attributes.Amount)
	if err != nil {
		Log(r).WithError(err).Error("Failed to update balance")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if !broadcastWithdrawalTx(req, r) {
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// see create_balance.go for explanation
	balance := getBalanceByDID(req.Data.ID, true, w, r)
	if balance == nil {
		return
	}

	ape.Render(w, newBalanceModel(*balance))
}

func isEnoughPoints(req resources.WithdrawRequest, w http.ResponseWriter, r *http.Request) bool {
	balance := getBalanceByDID(req.Data.ID, false, w, r)
	if balance == nil {
		return false
	}

	if balance.Amount < req.Data.Attributes.Amount {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"data/attributes/amount": fmt.Errorf("insufficient balance: %d", balance.Amount),
		})...)
		return false
	}

	return true
}

func broadcastWithdrawalTx(req resources.WithdrawRequest, r *http.Request) bool {
	var (
		from  = cosmos.MustAccAddressFromBech32(Broadcaster(r).Sender())
		to    = cosmos.MustAccAddressFromBech32(req.Data.Attributes.Address)
		urmo  = req.Data.Attributes.Amount * PointPrice(r)
		coins = cosmos.NewCoins(cosmos.NewInt64Coin("urmo", int64(urmo)))
	)

	err := Broadcaster(r).BroadcastTx(r.Context(), bank.NewMsgSend(from, to, coins))
	if err != nil {
		Log(r).WithError(err).Error("Failed to broadcast transaction")
		return false
	}

	return true
}
