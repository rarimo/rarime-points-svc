package handlers

import (
	"fmt"
	"net/http"

	"github.com/rarimo/decentralized-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func ClaimEvent(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewClaimEvent(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	event, err := EventsQ(r).FilterByID(req.Data.ID).FilterByStatus(data.EventFulfilled).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get event by balance ID")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if event == nil {
		Log(r).Debugf("Event not found for id=%s status=%s", req.Data.ID, data.EventFulfilled)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(event.Nullifier)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	evType := EventTypes(r).Get(event.Type, evtypes.FilterInactive)
	if evType == nil {
		Log(r).Infof("Event type %s is inactive", event.Type)
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	balance, err := BalancesQ(r).FilterByNullifier(event.Nullifier).FilterDisabled().Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by nullifier")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if balance == nil || balance.Country == nil {
		msg := "did not verify passport"
		if balance == nil {
			msg = "is disabled"
		}
		Log(r).Infof("Balance nullifier=%s %s", event.Nullifier, msg)
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	if !evType.IgnoreCountryLimit {
		country, err := CountriesQ(r).FilterByCodes(*balance.Country).Get()
		if err != nil || country == nil { // country must exist if no errors
			Log(r).WithError(err).Error("Failed to get country by code")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		if !country.ReserveAllowed {
			Log(r).Infof("Reserve is not allowed for country=%s", *balance.Country)
			ape.RenderErr(w, problems.Forbidden())
			return
		}
		if country.Reserved >= country.ReserveLimit {
			Log(r).Infof("Reserve limit is reached for country=%s", *balance.Country)
			ape.RenderErr(w, problems.Forbidden())
			return
		}
	}

	err = EventsQ(r).Transaction(func() error {
		event, err = claimEvent(r, event, balance)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		Log(r).WithError(err).Errorf("Failed to claim event %s and accrue %d points to the balance %s",
			event.ID, evType.Reward, event.Nullifier)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// balance should exist cause of previous logic
	balance, err = BalancesQ(r).GetWithRank(event.Nullifier)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by nullifier with rank")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newClaimEventResponse(*event, evType.Resource(), *balance))
}

// claimEvent requires event to exist
// call in transaction to prevent unexpected changes
func claimEvent(r *http.Request, event *data.Event, balance *data.Balance) (claimed *data.Event, err error) {
	if event == nil {
		return nil, nil
	}

	evType := EventTypes(r).Get(event.Type, evtypes.FilterInactive)
	if evType == nil {
		return event, nil
	}

	claimed, err = EventsQ(r).FilterByID(event.ID).Update(data.EventClaimed, nil, &evType.Reward)
	if err != nil {
		return nil, fmt.Errorf("update event status: %w", err)
	}

	err = DoClaimEventUpdates(
		Levels(r),
		ReferralsQ(r),
		BalancesQ(r),
		CountriesQ(r),
		*balance,
		evType.Reward,
		evType.IgnoreCountryLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to do claim event updates: %w", err)
	}

	return claimed, nil
}

// DoClaimEventUpdates do updates which link to claim event:
// update reserved amount in country;
// lvlup and update referrals count;
// accruing points;
//
// Balance must be active and with verified passport;
//
// ignoreCountryLimit determines whether to ignore checking and updating the country reserve.
// If true, the country reserve is not checked and updated.
func DoClaimEventUpdates(
	levels config.Levels,
	referralsQ data.ReferralsQ,
	balancesQ data.BalancesQ,
	countriesQ data.CountriesQ,
	balance data.Balance,
	reward int64,
	ignoreCountryLimit bool) (err error) {

	level, err := doLvlUpAndReferralsUpdate(levels, referralsQ, balance, reward)
	if err != nil {
		return fmt.Errorf("failed to do lvlup and referrals updates: %w", err)
	}

	err = balancesQ.FilterByNullifier(balance.Nullifier).Update(map[string]any{
		data.ColAmount: pg.AddToValue(data.ColAmount, reward),
		data.ColLevel:  level,
	})
	if err != nil {
		return fmt.Errorf("update balance amount and level: %w", err)
	}

	if !ignoreCountryLimit {
		err = countriesQ.FilterByCodes(*balance.Country).Update(map[string]any{
			data.ColReserved: pg.AddToValue(data.ColReserved, reward),
		})
		if err != nil {
			return fmt.Errorf("increase country reserve: %w", err)
		}
	}

	return nil
}

func doLvlUpAndReferralsUpdate(levels config.Levels, referralsQ data.ReferralsQ, balance data.Balance, reward int64) (level int, err error) {
	refsCount, level := levels.LvlUp(balance.Level, reward+balance.Amount)
	if refsCount > 0 {
		count, err := referralsQ.New().FilterByNullifier(balance.Nullifier).Count()
		if err != nil {
			return 0, fmt.Errorf("failed to get referral count: %w", err)
		}

		refToAdd := PrepareReferralsToAdd(balance.Nullifier, uint64(refsCount), count)
		if err = referralsQ.New().Insert(refToAdd...); err != nil {
			return 0, fmt.Errorf("failed to insert referrals: %w", err)
		}
	}

	return level, nil
}

func newClaimEventResponse(
	event data.Event,
	meta resources.EventStaticMeta,
	balance data.Balance,
) resources.EventResponse {

	eventModel := newEventModel(event, meta)
	eventModel.Relationships = &resources.EventRelationships{
		Balance: resources.Relation{
			Data: &resources.Key{
				ID:   balance.Nullifier,
				Type: resources.BALANCE,
			},
		},
	}

	resp := resources.EventResponse{Data: eventModel}
	inc := NewBalanceModel(balance)
	resp.Included.Add(&inc)

	return resp
}
