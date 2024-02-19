package handlers

import (
	"fmt"
	"net/http"

	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func VerifyPassport(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewVerifyPassport(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	balance, err := BalancesQ(r).FilterByDID(req.UserDID).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if balance == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	err = EventsQ(r).Transaction(func() error {
		// If you make this endpoint public, you should check the passport hash for
		// uniqueness and provide a better validation. Think about other changes too.
		err = BalancesQ(r).FilterByDID(req.UserDID).SetPassport(req.Hash, req.Expiry)
		if err != nil {
			return fmt.Errorf("set passport for balance by DID: %w", err)
		}

		evType := EventTypes(r).Get(evtypes.TypeReferralSpecific, evtypes.FilterInactive)
		if evType == nil {
			Log(r).Debug("Referral event type is disabled or expired, not accruing points to referrer")
			return nil
		}

		refDID, err := getReferrerDID(*balance, r)
		if err != nil {
			return fmt.Errorf("get referrer DID by referred_by: %w", err)
		}
		if refDID == "" {
			return nil
		}

		err = EventsQ(r).Insert(data.Event{
			UserDID: refDID,
			Type:    evType.Name,
			Status:  data.EventFulfilled,
		})
		if err != nil {
			return fmt.Errorf("add event for referrer: %w", err)
		}

		return nil
	})

	if err != nil {
		Log(r).WithError(err).Error("Failed to set passport and add event for referrer")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getReferrerDID(balance data.Balance, r *http.Request) (string, error) {
	if !balance.ReferredBy.Valid {
		return "", nil
	}

	refBy := balance.ReferredBy.String
	referrer, err := BalancesQ(r).FilterByReferralID(refBy).Get()
	if err != nil {
		return "", fmt.Errorf("failed to get balance by referral ID: %w", err)
	}
	if referrer == nil {
		return "", fmt.Errorf("referrer not found: %s", refBy)
	}

	Log(r).Debugf("Found referrer: DID=%s", referrer.DID)
	return referrer.DID, nil
}
