package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/rarimo/rarime-auth-svc/pkg/auth"
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

	event := getEventToClaim(req.Data.ID, w, r)
	if event == nil {
		return
	}
	balance := getBalanceByID(event.BalanceID, true, w, r)
	if balance == nil {
		return
	}

	evType := EventTypes(r).Get(event.Type)
	if evType == nil {
		Log(r).Error("Wrong event type is stored in DB: might be bad event config")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	event = claimEventWithPoints(*event, balance.Amount, int(evType.Reward), w, r)
	if event == nil {
		return
	}
	// can't return balance on update, see create_balance.go
	balance = getBalanceByID(event.BalanceID, true, w, r)
	if balance == nil {
		return
	}

	ape.Render(w, newClaimEventResponse(*event, *evType, *balance))
}

func getEventToClaim(id string, w http.ResponseWriter, r *http.Request) *data.Event {
	event, err := EventsQ(r).
		FilterByID(id).
		FilterByStatus(data.EventFulfilled).
		Get()

	if err != nil {
		Log(r).WithError(err).Error("Failed to get event by balance ID")
		ape.RenderErr(w, problems.InternalError())
		return nil
	}

	if event == nil {
		Log(r).Debugf("Event not found for id=%s status=%s", id, data.EventFulfilled)
		ape.RenderErr(w, problems.NotFound())
		return nil
	}

	return event
}

func getBalanceByID(id string, doAuth bool, w http.ResponseWriter, r *http.Request) *data.Balance {
	balance, err := BalancesQ(r).WithRank().FilterByID(id).Get()

	if err != nil || balance == nil {
		if err == nil {
			err = fmt.Errorf("DB constraint violation: found event with balance_id=%s not present", id)
		}

		Log(r).WithError(err).Error("Failed to get balance by ID")
		ape.RenderErr(w, problems.InternalError())
		return nil
	}

	if doAuth && !auth.Authenticates(UserClaims(r), auth.UserGrant(balance.DID)) {
		ape.RenderErr(w, problems.Unauthorized())
		return nil
	}

	return balance
}

func claimEventWithPoints(event data.Event, currBalance, reward int, w http.ResponseWriter, r *http.Request) *data.Event {
	claimed := data.Event{
		ID:     event.ID,
		Status: data.EventClaimed,
		PointsAmount: sql.NullInt32{
			Int32: int32(reward),
			Valid: true,
		},
	}

	err := EventsQ(r).Update(claimed)
	if err != nil {
		Log(r).WithError(err).Error("Failed to claim event")
		ape.RenderErr(w, problems.InternalError())
		return nil
	}

	err = BalancesQ(r).FilterByID(event.BalanceID).UpdateAmount(currBalance + reward)
	if err != nil {
		Log(r).WithError(err).Error("Failed to accrue points to the balance")
		ape.RenderErr(w, problems.InternalError())
		return nil
	}
	// While we don't have updated_at and other special attributes in events, we can
	// safely return the same struct without redundant queries. It is still faster
	// than with RETURNING clause.
	event.Status = claimed.Status
	event.PointsAmount = claimed.PointsAmount
	return &event
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
				ID:   balance.ID,
				Type: resources.BALANCE,
			},
		},
	}

	resp := resources.EventResponse{Data: eventModel}
	inc := newBalanceModel(balance)
	resp.Included.Add(&inc)

	return resp
}
