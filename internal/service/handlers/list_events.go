package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rarimo/points-svc/internal/data"
	"github.com/rarimo/points-svc/internal/service/requests"
	"github.com/rarimo/points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func ListEvents(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewListEvents(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	balance := getBalanceByDID(req.DID, false, w, r)
	if balance == nil {
		return
	}

	q := EventsQ(r).FilterByBalanceID(balance.ID)
	if req.FilterStatus != nil {
		q.FilterByStatus(*req.FilterStatus)
	}

	events, err := q.Page(&req.CursorPageParams).Select()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get event list")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var last string
	if len(events) > 0 {
		last = events[len(events)-1].ID
	}

	resp := newEventsResponse(events)
	resp.Links = req.CursorParams.GetCursorLinks(r, last)
	ape.Render(w, resp)
}

func newEventsResponse(events []data.Event) *resources.EventListResponse {
	list := make([]resources.Event, len(events))

	for i, event := range events {
		var dynamic *json.RawMessage
		if event.Meta.Valid {
			d := json.RawMessage(event.Meta.String)
			dynamic = &d
		}

		list[i] = resources.Event{
			Key: resources.Key{
				ID:   event.ID,
				Type: resources.EVENT,
			},
			Attributes: resources.EventAttributes{
				CreatedAt: event.CreatedAt,
				Meta: resources.EventMeta{
					Static:  json.RawMessage{}, // TODO: add from config
					Dynamic: dynamic,
				},
				Status: event.Status.String(),
			},
		}
	}

	return &resources.EventListResponse{Data: list}
}
