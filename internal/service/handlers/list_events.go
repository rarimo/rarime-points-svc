package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rarimo/auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func ListEvents(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewListEvents(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(req.FilterDID)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	events, err := EventsQ(r).
		FilterByUserDID(req.FilterDID).
		FilterByStatus(req.FilterStatus...).
		FilterByType(req.FilterType...).
		Page(&req.CursorPageParams).
		Select()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get event list")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var eventsCount int
	if req.Count {
		eventsCount, err = EventsQ(r).
			FilterByUserDID(req.FilterDID).
			FilterByStatus(req.FilterStatus...).
			FilterByType(req.FilterType...).
			Count()
		if err != nil {
			Log(r).WithError(err).Error("Failed to count events")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	meta, err := getOrderedEventsMeta(events, r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get ordered events metadata")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var last string
	if len(events) > 0 {
		last = events[len(events)-1].ID
	}

	resp := newEventsResponse(events, meta)
	resp.Links = req.CursorParams.GetLinks(r, last)
	if req.Count {
		_ = resp.PutMeta(struct {
			EventsCount int `json:"events_count"`
		}{eventsCount})
	}
	ape.Render(w, resp)
}

func getOrderedEventsMeta(events []data.Event, r *http.Request) ([]resources.EventStaticMeta, error) {
	res := make([]resources.EventStaticMeta, len(events))

	for i, event := range events {
		// even if event type was disabled, we should return it from history
		evType := EventTypes(r).Get(event.Type)
		if evType == nil {
			return nil, errors.New("wrong event type is stored in DB: might be bad event config")
		}
		res[i] = evType.Resource()
	}

	return res, nil
}

func newEventModel(event data.Event, meta resources.EventStaticMeta) resources.Event {
	return resources.Event{
		Key: resources.Key{
			ID:   event.ID,
			Type: resources.EVENT,
		},
		Attributes: resources.EventAttributes{
			CreatedAt: event.CreatedAt,
			UpdatedAt: event.UpdatedAt,
			Meta: resources.EventMeta{
				Static:  meta,
				Dynamic: (*json.RawMessage)(&event.Meta),
			},
			Status:       event.Status.String(),
			PointsAmount: event.PointsAmount,
		},
	}
}

func newEventsResponse(events []data.Event, meta []resources.EventStaticMeta) *resources.EventListResponse {
	list := make([]resources.Event, len(events))
	for i, event := range events {
		list[i] = newEventModel(event, meta[i])
	}

	return &resources.EventListResponse{Data: list}
}
