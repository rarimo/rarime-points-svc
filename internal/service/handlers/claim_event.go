package handlers

import (
	"cosmossdk.io/errors"
	"fmt"
	"net/http"

	"github.com/rarimo/auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data"
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

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(event.UserDID)) {
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

	balance, err := BalancesQ(r).FilterByDID(event.UserDID).FilterDisabled().Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if balance == nil {
		Log(r).Infof("Attempt to claim: balance user_did=%s is disabled", event.UserDID)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	event, err = claimEventWithPoints(*event, evType.Reward, r)
	if err != nil {
		Log(r).WithError(err).Errorf("Failed to claim event %s and accrue %d points to the balance %s",
			event.ID, evType.Reward, event.UserDID)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if err := updateLevelAfterBalance(event.UserDID, r); err != nil {
		Log(r).WithError(err).Errorf("Failed to update level for user %s", event.UserDID)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// balance should exist cause of previous logic
	balance, err = BalancesQ(r).GetWithRank(event.UserDID)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID with rank")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newClaimEventResponse(*event, evType.Resource(), *balance))
}

// requires: event exist
func claimEventWithPoints(event data.Event, reward int64, r *http.Request) (claimed *data.Event, err error) {
	err = EventsQ(r).Transaction(func() error {
		updated, err := EventsQ(r).FilterByID(event.ID).Update(data.EventClaimed, nil, &reward)
		if err != nil {
			return fmt.Errorf("update event status: %w", err)
		}

		err = BalancesQ(r).FilterByDID(event.UserDID).UpdateAmountBy(reward)
		if err != nil {
			return fmt.Errorf("update balance amount: %w", err)
		}

		claimed = updated
		return nil
	})
	return
}

func updateLevelAfterBalance(did string, r *http.Request) error {
	return BalancesQ(r).Transaction(func() error {
		balance, err := BalancesQ(r).FilterByDID(did).Get()
		if err != nil {
			return err
		}

		level := getLevelByAmount(balance.Amount, r)
		claimId := balance.LevelClaimId

		if level > balance.Level {
			if claimId != nil {
				if err := Issuer(r).RevokeClaim(*claimId); err != nil {
					return errors.Wrap(err, "failed to revoke claim")
				}
			}

			claimId, err = Issuer(r).IssueLevelClaim(balance.DID, level)
			if err != nil {
				return errors.Wrap(err, "failed to issuer claim")
			}
		}

		if level != balance.Level {
			return BalancesQ(r).FilterByDID(did).SetLevel(level, *claimId)
		}
		return nil
	})
}

func getLevelByAmount(amount int64, r *http.Request) int {
	for i, am := range Levels(r) {
		if amount < am {
			return i
		}
	}

	return len(Levels(r))
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
				ID:   balance.DID,
				Type: resources.BALANCE,
			},
		},
	}

	resp := resources.EventResponse{Data: eventModel}
	inc := newBalanceModel(balance)
	resp.Included.Add(&inc)

	return resp
}
