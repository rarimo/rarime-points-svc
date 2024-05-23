package handlers

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/rarimo/decentralized-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
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

	evType := EventTypes(r).Get(event.Type) // expired events can be claimed
	if evType == nil {
		Log(r).Errorf("Wrong event type %s is stored in DB: might be bad event config", event.Type)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if evType.Disabled {
		Log(r).Infof("Attempt to claim: event type %s is disabled", event.Type)
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	balance, err := BalancesQ(r).FilterByNullifier(event.Nullifier).FilterDisabled().Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by nullifier")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if balance == nil {
		Log(r).Infof("Attempt to claim: balance nullifier=%s is disabled", event.Nullifier)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	event, err = claimEventWithPoints(r, *event, evType.Reward, balance)
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

// requires: event exist
func claimEventWithPoints(r *http.Request, event data.Event, reward int64, balance *data.Balance) (claimed *data.Event, err error) {
	err = EventsQ(r).Transaction(func() error {
		refsCount, level := levelsToUp(r, balance.Level, reward+balance.Amount)
		if level != 0 {
			count, err := ReferralsQ(r).FilterByNullifier(event.Nullifier).Count()
			if err != nil {
				return errors.Wrap(err, "failed to get referral count")
			}

			refToAdd := prepareReferralsToAdd(event.Nullifier, uint64(refsCount), count)
			if err = ReferralsQ(r).Insert(refToAdd...); err != nil {
				return errors.Wrap(err, "failed to insert referrals")
			}

			if err = BalancesQ(r).FilterByNullifier(event.Nullifier).SetLevel(level); err != nil {
				return errors.Wrap(err, "failed to update level")
			}
		}

		updated, err := EventsQ(r).FilterByID(event.ID).Update(data.EventClaimed, nil, &reward)
		if err != nil {
			return fmt.Errorf("update event status: %w", err)
		}

		err = BalancesQ(r).FilterByNullifier(event.Nullifier).UpdateAmountBy(reward)
		if err != nil {
			return fmt.Errorf("update balance amount: %w", err)
		}

		claimed = updated
		return nil
	})
	return
}

// Take summary points after claim and calculate new lvls
func levelsToUp(r *http.Request, level int, points int64) (referralsToAdd int, newLevel int) {
	lvls := make([]int, 0, len(Levels(r)))
	for k, v := range Levels(r) {
		if k <= level {
			continue
		}
		if int64(v.Threshold) > points {
			break
		}

		referralsToAdd += v.Referrals
		lvls = append(lvls, k)
	}

	newLevel = slices.Max(lvls)
	return
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
	inc := newBalanceModel(balance)
	resp.Included.Add(&inc)

	return resp
}
