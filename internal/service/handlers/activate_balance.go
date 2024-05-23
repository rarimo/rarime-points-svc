package handlers

import (
	"fmt"
	"net/http"

	"github.com/rarimo/decentralized-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func ActivateBalance(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewActivateBalance(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	nullifier := req.Data.ID

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(nullifier)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	balance, err := BalancesQ(r).FilterByNullifier(nullifier).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by nullifier")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if balance == nil {
		Log(r).Debug("Balance not exist")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	// Balance should be inactive
	if balance.ReferredBy.Valid {
		Log(r).Debug("Balance already activated")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	referral, err := ReferralsQ(r).FilterByIsConsumed(false).Get(req.Data.Attributes.ReferredBy)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get referral by ID")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if referral == nil {
		Log(r).Debug("Referral code already consumed or not exists")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	referrals := prepareReferralsToAdd(nullifier, 5, 0)

	err = EventsQ(r).Transaction(func() error {
		Log(r).Debugf("%s referral code will be set for nullifier=%s", req.Data.Attributes.ReferredBy, nullifier)
		if err = BalancesQ(r).FilterByNullifier(nullifier).SetReferredBy(req.Data.Attributes.ReferredBy); err != nil {
			return fmt.Errorf("set referred_by: %w", err)
		}

		Log(r).Debugf("%d referrals will be added for nullifier=%s", len(referrals), nullifier)
		if err = ReferralsQ(r).Insert(referrals...); err != nil {
			return fmt.Errorf("add referrals: %w", err)
		}

		Log(r).Debugf("%s referral will be consumed", req.Data.Attributes.ReferredBy)
		if _, err = ReferralsQ(r).Consume(req.Data.Attributes.ReferredBy); err != nil {
			return fmt.Errorf("consume referral: %w", err)
		}

		if balance.PassportHash.Valid {
			evType := EventTypes(r).Get(evtypes.TypeReferralSpecific, evtypes.FilterInactive)
			if evType == nil {
				Log(r).Debug("Referral event type is disabled or expired, not accruing points to referrer")
				return nil
			}

			err = EventsQ(r).Insert(data.Event{
				Nullifier: referral.Nullifier,
				Type:      evType.Name,
				Status:    data.EventFulfilled,
				Meta:      data.Jsonb(fmt.Sprintf(`{"nullifier": "%s"}`, nullifier)),
			})
			if err != nil {
				return fmt.Errorf("add event for referrer: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		Log(r).WithError(err).Error("Failed to activate balance")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	balance, err = BalancesQ(r).GetWithRank(nullifier)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by nullifier with rank")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newBalanceResponse(*balance, referrals))
}
