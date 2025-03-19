package handlers

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/service/referralid"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func EditReferrals(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewEditReferrals(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	balance, err := BalancesQ(r).FilterByNullifier(req.Nullifier).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by nullifier")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if balance == nil {
		if req.Count == 0 {
			Log(r).Debugf("Balance %s not found, skipping creation for count=0", req.Nullifier)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var code string
		err = EventsQ(r).Transaction(func() error {
			events := prepareEventsWithRef(req.Nullifier, "", r)
			if err = createBalanceWithEvents(req.Nullifier, "", events, r); err != nil {
				return fmt.Errorf("failed to create balance with events: %w", err)
			}

			code = referralid.New(req.Nullifier, 0)
			err = ReferralsQ(r).Insert(data.Referral{
				ID:        code,
				Nullifier: req.Nullifier,
				UsageLeft: int32(req.Count),
			})
			if err != nil {
				return fmt.Errorf("failed to insert referral for nullifier [%s]: %w", req.Nullifier, err)
			}

			return nil
		})
		if err != nil {
			Log(r).WithError(err).Errorf("failed to create genesis balance [%s]", req.Nullifier)
			ape.RenderErr(w, problems.InternalError())
			return
		}

		ape.Render(w, struct {
			Ref       string `json:"referral"`
			UsageLeft uint64 `json:"usage_left"`
		}{code, req.Count})
		return
	}

	if balance.ReferredBy.Valid {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"balance": fmt.Errorf("genesis balances must be inactive")})...)
		return
	}

	referrals, err := ReferralsQ(r).FilterByNullifier(req.Nullifier).Select()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to select referrals for nullifier [%s]", req.Nullifier)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if len(referrals) != 1 {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"balance": fmt.Errorf("genesis balances must have only one referral")})...)
		return
	}

	referral, err := ReferralsQ(r).FilterByNullifier(req.Nullifier).Update(int(req.Count))
	if err != nil {
		Log(r).WithError(err).Errorf("failed to update referral usage count for nullifier [%s]", req.Nullifier)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if referral == nil {
		Log(r).Errorf("critical: referral absent for user [%s]", req.Nullifier)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, struct {
		Ref       string `json:"referral"`
		UsageLeft uint64 `json:"usage_left"`
	}{
		referral.ID,
		uint64(referral.UsageLeft),
	})

}

func PrepareReferralsToAdd(nullifier string, count, index uint64) []data.Referral {
	refCodes := referralid.NewMany(nullifier, count, index)
	refs := make([]data.Referral, len(refCodes))

	for i, code := range refCodes {
		refs[i] = data.Referral{
			ID:        code,
			Nullifier: nullifier,
			UsageLeft: 1,
		}
	}

	return refs
}
