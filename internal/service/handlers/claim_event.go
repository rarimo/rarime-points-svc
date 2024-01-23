package handlers

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func ClaimEvent(w http.ResponseWriter, r *http.Request) {
	did := r.Header.Get("X-User-DID")

	eventID := chi.URLParam(r, "id")
	if eventID == "" {
		ape.RenderErr(w, problems.BadRequest(nil)...)
		return
	}

	balance := getBalanceByDID(did, false, w, r)
	if balance == nil {
		return
	}

	event, err := EventsQ(r).
		FilterByID(eventID).
		FilterByBalanceID(balance.ID).
		FilterByStatus(data.EventFulfilled).
		Get()

	if err != nil {
		Log(r).WithError(err).Error("Failed to get event by balance ID")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if event == nil {
		Log(r).Debugf("Event not found for id=%s balance_id=%s status=%s", eventID, balance.ID, data.EventFulfilled)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	evType := EventTypes(r).Get(event.Type)
	if evType == nil {
		Log(r).Error("Wrong event type is stored in DB: might be bad event config")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	err = EventsQ(r).Update(data.Event{
		ID:     event.ID,
		Status: data.EventClaimed,
		PointsAmount: sql.NullInt32{
			Int32: evType.Reward,
			Valid: true,
		},
	})
	if err != nil {
		Log(r).WithError(err).Error("Failed to claim event")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	err = BalancesQ(r).FilterByID(balance.ID).UpdateAmount(int(evType.Reward))
	if err != nil {
		Log(r).WithError(err).Error("Failed to accrue points to the balance")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
