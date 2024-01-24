package handlers

import (
	"encoding/json"
	"net/http"

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

	balance := getBalanceByDID(req.FilterDID, false, w, r)
	if balance == nil {
		return
	}

	events, err := EventsQ(r).
		FilterByBalanceID(balance.ID).
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
			FilterByBalanceID(balance.ID).
			FilterByStatus(req.FilterStatus...).
			FilterByType(req.FilterType...).
			Count()
		if err != nil {
			Log(r).WithError(err).Error("Failed to count events")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	meta, ok := getOrderedEventsMeta(events, w, r)
	if !ok {
		return
	}

	var last string
	if len(events) > 0 {
		last = events[len(events)-1].ID
	}

	resp := newEventsResponse(events, meta)
	resp.Links = req.CursorParams.GetCursorLinks(r, last)
	if req.Count {
		_ = resp.PutMeta(struct {
			EventsCount int `json:"events_count"`
		}{eventsCount})
	}
	ape.Render(w, resp)
}

func getOrderedEventsMeta(events []data.Event, w http.ResponseWriter, r *http.Request) ([]resources.EventStaticMeta, bool) {
	res := make([]resources.EventStaticMeta, len(events))

	for i, event := range events {
		evType := EventTypes(r).Get(event.Type)
		if evType == nil {
			Log(r).Error("Wrong event type is stored in DB: might be bad event config")
			ape.RenderErr(w, problems.InternalError())
			return nil, false
		}
		res[i] = *evType
	}

	return res, true
}

func newEventModel(event data.Event, meta resources.EventStaticMeta) resources.Event {
	var dynamic *json.RawMessage
	if event.Meta.Valid {
		d := json.RawMessage(event.Meta.String)
		dynamic = &d
	}

	var points *int32
	if event.PointsAmount.Valid {
		points = &event.PointsAmount.Int32
	}

	return resources.Event{
		Key: resources.Key{
			ID:   event.ID,
			Type: resources.EVENT,
		},
		Attributes: resources.EventAttributes{
			CreatedAt: event.CreatedAt,
			Meta: resources.EventMeta{
				Static:  meta,
				Dynamic: dynamic,
			},
			Status:       event.Status.String(),
			PointsAmount: points,
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
