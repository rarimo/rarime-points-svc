package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	api "github.com/rarimo/rarime-points-svc/pkg/connector"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func FulfillEvent(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewFulfillEvent(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if apiErr := validateEventType(req.EventType, r); apiErr != nil {
		ape.RenderErr(w, apiErr)
		return
	}

	internalErr := api.CodeInternalError.JSONAPIError()
	balance, err := BalancesQ(r).FilterByDID(req.UserDID).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, internalErr)
		return
	}

	if balance == nil {
		if req.ExternalID != nil {
			ape.RenderErr(w, api.CodeDidUnknown.JSONAPIError())
			return
		}
		// it would be more efficient to insert already fulfilled event with the specific type
		if err = createBalanceWithEvents(req.UserDID, "", r); err != nil {
			Log(r).WithError(err).Error("Failed to create balance with events")
			ape.RenderErr(w, internalErr)
			return
		}
	}

	eventID, err := getEventToFulfill(req, r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get event to fulfill")
		ape.RenderErr(w, internalErr)
		return
	}
	if eventID == "" {
		ape.RenderErr(w, api.CodeEventNotFound.JSONAPIError())
		return
	}

	_, err = EventsQ(r).FilterByID(eventID).Update(data.EventFulfilled, nil, nil)
	if err != nil {
		Log(r).WithError(err).Error("Failed to update event")
		ape.RenderErr(w, internalErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validateEventType(name string, r *http.Request) *jsonapi.ErrorObject {
	evType := EventTypes(r).Get(name)

	switch {
	case evType == nil || evType.Disabled:
		return api.CodeEventDisabled.JSONAPIError()
	case evtypes.FilterExpired(*evType):
		return api.CodeEventExpired.JSONAPIError()
	}

	return nil
}

func getEventToFulfill(req api.FulfillEventRequest, r *http.Request) (eventID string, err error) {
	q := EventsQ(r).FilterByUserDID(req.UserDID).
		FilterByType(req.EventType).
		FilterByStatus(data.EventOpen)
	if req.ExternalID != nil {
		q.FilterByExternalID(*req.ExternalID)
	}

	events, err := q.Select()
	switch {
	case err != nil:
		return "", fmt.Errorf("select events: %w", err)
	case len(events) > 1:
		return "", fmt.Errorf("multiple events to fulfill found: %d", len(events))
	case len(events) == 0:
		return "", nil
	}

	return events[0].ID, nil
}
