package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/rarimo/auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/referralid"
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

	balance, err := getBalanceByDID(did, false, r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// Balance should not exist
	if balance != nil {
		ape.RenderErr(w, problems.Conflict())
		return
	}

	var referredBy string
	if attr := req.Data.Attributes; attr != nil {
		referrer, err := BalancesQ(r).FilterByReferralID(attr.ReferredBy).Get()
		if err != nil {
			Log(r).WithError(err).Error("Failed to check referrer existence")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		if referrer == nil {
			Log(r).Debugf("Referrer not found for referral_id=%s", attr.ReferredBy)
			ape.RenderErr(w, problems.NotFound())
			return
		}
		referredBy = attr.ReferredBy
	}

	if err = createBalanceWithEvents(did, referredBy, r); err != nil {
		Log(r).WithError(err).Error("Failed to create balance with events")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// We can't return inserted balance in a single query, because we can't calculate
	// rank in transaction: RANK() is a window function allowed on a set of rows,
	// while with RETURNING we operate a single one.
	// Balance will exist cause of previous logic.
	balance, err = getBalanceByDID(did, true, r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newBalanceModel(*balance))
}

func createBalanceWithEvents(did, refBy string, r *http.Request) error {
	return EventsQ(r).Transaction(func() error {
		err := BalancesQ(r).Insert(data.Balance{
			DID:        did,
			ReferralID: referralid.New(did),
			ReferredBy: sql.NullString{String: refBy, Valid: refBy != ""},
		})

		if err != nil {
			return fmt.Errorf("add balance: %w", err)
		}

		events := prepareEventsWithRef(did, refBy, r)
		Log(r).Debugf("%d events will be added for user_did=%s", len(events), did)

		if err = EventsQ(r).Insert(events...); err != nil {
			return fmt.Errorf("add open events: %w", err)
		}

		return nil
	})
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
