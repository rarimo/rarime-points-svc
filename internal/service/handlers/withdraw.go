package handlers

import (
	"fmt"
	"net/http"
	"time"

	cosmos "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

	if !isEligibleToWithdraw(req, w, r) {
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

	// see create_balance.go for explanation
	balance := getBalanceByDID(req.Data.ID, true, w, r)
	if balance == nil {
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

func isEligibleToWithdraw(req resources.WithdrawRequest, w http.ResponseWriter, r *http.Request) bool {
	balance := getBalanceByDID(req.Data.ID, false, w, r)
	if balance == nil {
		return false
	}

	render := func(field, format string, a ...any) bool {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			field: fmt.Errorf(format, a...),
		})...)
		return false
	}

	if !balance.PassportHash.Valid {
		return render("is_verified", "user must have verified passport for withdrawals")
	}
	if balance.PassportExpires.Time.Before(time.Now().UTC()) {
		return render("is_verified", "user passport is expired")
	}
	if balance.Amount < req.Data.Attributes.Amount {
		return render("data/attributes/amount", "insufficient balance: %d", balance.Amount)
	}

	return true
}

func broadcastWithdrawalTx(req resources.WithdrawRequest, r *http.Request) error {
	var (
		from  = cosmos.MustAccAddressFromBech32(Broadcaster(r).Sender())
		to    = cosmos.MustAccAddressFromBech32(req.Data.Attributes.Address)
		urmo  = req.Data.Attributes.Amount * PointPrice(r)
		coins = cosmos.NewCoins(cosmos.NewInt64Coin("urmo", int64(urmo)))
	)

	err := Broadcaster(r).BroadcastTx(r.Context(), bank.NewMsgSend(from, to, coins))
	if err != nil {
		return fmt.Errorf("broadcast withdrawal tx: %w", err)
	}

	return nil
}
