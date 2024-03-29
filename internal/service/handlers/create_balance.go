package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/rarimo/auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateBalance(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewCreateBalance(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	did := req.Data.ID

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(did)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	balance, err := BalancesQ(r).FilterByDID(did).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// Balance should not exist or be inactive
	if balance != nil && balance.ReferredBy.Valid {
		ape.RenderErr(w, problems.Conflict())
		return
	}

	referral, err := ReferralsQ(r).FilterByIsConsumed(false).Get(req.Data.Attributes.ReferredBy)
	if referral == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if err != nil {
		Log(r).WithError(err).Error("Failed to get referral by ID")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	referrals := prepareReferralsToAdd(did, 5, 0)

	if balance == nil {
		events := prepareEventsWithRef(did, req.Data.Attributes.ReferredBy, r)
		if err = createBalanceWithEventsAndReferrals(did, req.Data.Attributes.ReferredBy, events, referrals, r); err != nil {
			Log(r).WithError(err).Error("Failed to create balance with events")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		// We can't return inserted balance in a single query, because we can't calculate
		// rank in transaction: RANK() is a window function allowed on a set of rows,
		// while with RETURNING we operate a single one.
		// Balance will exist cause of previous logic.
		balance, err = BalancesQ(r).GetWithRank(did)
		if err != nil {
			Log(r).WithError(err).Error("Failed to get created balance by DID")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		ape.Render(w, newBalanceResponse(*balance, referrals))
		return
	}

	err = EventsQ(r).Transaction(func() error {
		Log(r).Debugf("%s referral code will be added for user_did=%s", req.Data.Attributes.ReferredBy, did)
		if err = BalancesQ(r).FilterByDID(did).SetReferredBy(req.Data.Attributes.ReferredBy); err != nil {
			return fmt.Errorf("set referred_by: %w", err)
		}

		Log(r).Debugf("%d referrals will be added for user_did=%s", len(referrals), did)
		if err = ReferralsQ(r).Insert(referrals...); err != nil {
			return fmt.Errorf("add referrals: %w", err)
		}

		Log(r).Debugf("%s referral will be consumed", req.Data.Attributes.ReferredBy)
		if _, err = ReferralsQ(r).Consume(req.Data.Attributes.ReferredBy); err != nil {
			return fmt.Errorf("consume referral: %w", err)
		}
		return nil
	})

	if err != nil {
		Log(r).WithError(err).Error("Failed to activate balance")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	balance, err = BalancesQ(r).GetWithRank(did)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID with rank")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newBalanceResponse(*balance, referrals))
}

func prepareEventsWithRef(did, refBy string, r *http.Request) []data.Event {
	events := EventTypes(r).PrepareEvents(did, evtypes.FilterNotOpenable)
	if refBy == "" {
		return events
	}

	refType := EventTypes(r).Get(evtypes.TypeBeReferred, evtypes.FilterInactive)
	if refType == nil {
		Log(r).Debug("Referral event is disabled or expired, skipping it")
		return events
	}

	Log(r).WithFields(map[string]any{"user_did": did, "referred_by": refBy}).
		Debug("`Be referred` event will be added for referee user")

	return append(events, data.Event{
		UserDID: did,
		Type:    evtypes.TypeBeReferred,
		Status:  data.EventFulfilled,
	})
}

func createBalanceWithEvents(did, refBy string, events []data.Event, r *http.Request) error {
	return EventsQ(r).Transaction(func() error {
		err := BalancesQ(r).Insert(data.Balance{
			DID:        did,
			ReferredBy: sql.NullString{String: refBy, Valid: refBy != ""},
		})

		if err != nil {
			return fmt.Errorf("add balance: %w", err)
		}

		Log(r).Debugf("%d events will be added for user_did=%s", len(events), did)
		if err = EventsQ(r).Insert(events...); err != nil {
			return fmt.Errorf("add open events: %w", err)
		}

		return nil
	})
}

func createBalanceWithEventsAndReferrals(did, refBy string, events []data.Event, refCodes []data.Referral, r *http.Request) error {
	return EventsQ(r).Transaction(func() error {
		err := BalancesQ(r).Insert(data.Balance{
			DID:        did,
			ReferredBy: sql.NullString{String: refBy, Valid: refBy != ""},
		})

		if err != nil {
			return fmt.Errorf("add balance: %w", err)
		}

		Log(r).Debugf("%d events will be added for user_did=%s", len(events), did)
		if err = EventsQ(r).Insert(events...); err != nil {
			return fmt.Errorf("add open events: %w", err)
		}

		Log(r).Debugf("%d referrals will be added for user_did=%s", len(refCodes), did)
		if err = ReferralsQ(r).Insert(refCodes...); err != nil {
			return fmt.Errorf("add referrals: %w", err)
		}

		Log(r).Debugf("%s referral will be consumed", refBy)
		if _, err = ReferralsQ(r).Consume(refBy); err != nil {
			return fmt.Errorf("consume referral: %w", err)
		}

		return nil
	})
}
