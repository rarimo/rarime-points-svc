//nolint:all
package handlers

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rarimo/rarime-points-svc/internal/service/broadcaster"
	"github.com/rarimo/rarime-points-svc/internal/service/handlers"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/jsonapi"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
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
	log := handlers.Log(r).WithFields(map[string]any{
		"nullifier":     req.Data.ID,
		"points_amount": req.Data.Attributes.Amount,
		"dest_address":  req.Data.Attributes.Address,
	})

	if handlers.PointPrice(r).Disabled {
		log.Debug("Withdrawal is disabled")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	var (
		nullifier = req.Data.ID
		proof     = req.Data.Attributes.Proof
	)

	balance, errs := getAndVerifyBalanceEligibility(r, nullifier, nil)
	if len(errs) > 0 {
		ape.RenderErr(w, errs...)
		return
	}

	// validated in requests.NewWithdraw
	// addr, _ := cosmos.AccAddressFromBech32(req.Data.Attributes.Address)
	// never panics because of request validation
	// proof.PubSignals[zk.Nullifier] = mustHexToInt(nullifier)

	// err = Verifier(r).VerifyProof(proof, zk.WithEventData(addr))
	// if err != nil {
	// 	ape.RenderErr(w, problems.BadRequest(err)...)
	// 	return
	// }

	countryCode, err := requests.ExtractCountry(proof)
	if err != nil {
		log.WithError(err).Error("Critical: invalid country code provided, while the proof was valid")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	country, err := getOrCreateCountry(handlers.CountriesQ(r), countryCode) // +1 query is not critical
	if err != nil {
		log.WithError(err).Error("Failed to get or create country")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	errs = isEligibleToWithdraw(r, balance, req.Data.Attributes.Amount, *country)
	if len(errs) > 0 {
		ape.RenderErr(w, errs...)
		return
	}

	var withdrawal *data.Withdrawal
	err = handlers.EventsQ(r).Transaction(func() error {
		err = handlers.BalancesQ(r).FilterByNullifier(nullifier).Update(map[string]any{
			data.ColAmount: pg.AddToValue(data.ColAmount, -req.Data.Attributes.Amount),
		})
		if err != nil {
			return fmt.Errorf("decrease points amount: %w", err)
		}

		err = handlers.CountriesQ(r).FilterByCodes(*balance.Country).Update(map[string]any{
			data.ColWithdrawn: pg.AddToValue(data.ColWithdrawn, req.Data.Attributes.Amount),
		})
		if err != nil {
			return fmt.Errorf("increase country withdrawn: %w", err)
		}

		withdrawal, err = handlers.WithdrawalsQ(r).Insert(data.Withdrawal{
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
	balance, err = handlers.BalancesQ(r).GetWithRank(nullifier)
	if err != nil {
		log.WithError(err).Error("Failed to get balance by nullifier with rank")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newWithdrawResponse(*withdrawal, *balance))
}

func newWithdrawResponse(w data.Withdrawal, balance data.Balance) *resources.WithdrawalResponse {
	wm := handlers.NewWithdrawalModel(w)
	wm.Relationships = &resources.WithdrawalRelationships{
		Balance: resources.Relation{
			Data: &resources.Key{
				ID:   balance.Nullifier,
				Type: resources.BALANCE,
			},
		},
	}

	resp := resources.WithdrawalResponse{Data: wm}
	bm := handlers.NewBalanceModel(balance)
	resp.Included.Add(&bm)

	return &resp
}

func isEligibleToWithdraw(
	r *http.Request,
	balance *data.Balance,
	amount int64,
	country data.Country,
) []*jsonapi.ErrorObject {

	mapValidationErr := func(field, format string, a ...any) []*jsonapi.ErrorObject {
		return problems.BadRequest(validation.Errors{
			field: fmt.Errorf(format, a...),
		})
	}

	switch {
	case !balance.IsPassportProven:
		return mapValidationErr("data/attributes/proof", "passport must be proven beforehand")
	case balance.Amount < amount:
		return mapValidationErr("data/attributes/amount", "insufficient balance: %d", balance.Amount)
	case !country.WithdrawalAllowed:
		return mapValidationErr("country", "withdrawal is not allowed for country=%s", country.Code)
	case !handlers.Levels(r)[balance.Level].WithdrawalAllowed:
		return mapValidationErr("level", "must up level to have withdraw ability")
	case balance.Country != nil && *balance.Country != country.Code:
		return mapValidationErr("country", "country mismatch in proof and balance: %s", *balance.Country)
	}

	return nil
}

func broadcastWithdrawalTx(req resources.WithdrawRequest, r *http.Request) error {
	b, err := broadcaster.New(handlers.Broadcaster(r), handlers.Log(r))
	if err != nil {
		return fmt.Errorf("failed to create broadcaster: %w", err)
	}

	toAddr := common.HexToAddress(req.Data.Attributes.Address)
	amount := big.NewInt(req.Data.Attributes.Amount * handlers.PointPrice(r).PointPriceURMO)

	if err := b.BroadcastTransfer(r.Context(), toAddr, amount); err != nil {
		return fmt.Errorf("failed to broadcast ERC20 transfer: %w", err)
	}

	return nil
}
