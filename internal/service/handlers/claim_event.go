package handlers

import (
	"fmt"
	"net/http"

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

	evType := EventTypes(r).Get(event.Type)
	if evType == nil {
		Log(r).Error("Wrong event type is stored in DB: might be bad event config")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	event = claimEventWithPoints(*event, evType.Reward, w, r)
	if event == nil {
		return
	}
	// can't return balance with rank on update, see create_balance.go
	balance := getBalanceByDID(event.UserDID, true, w, r)
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

func claimEventWithPoints(event data.Event, reward int32, w http.ResponseWriter, r *http.Request) (claimed *data.Event) {
	err := EventsQ(r).Transaction(func() error {
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

	if err != nil {
		Log(r).WithError(err).Error("Failed to claim event and accrue points to the balance")
		ape.RenderErr(w, problems.InternalError())
	}

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
