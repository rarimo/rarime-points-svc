package handlers

import (
	"database/sql"
	"net/http"

	"github.com/rarimo/points-svc/internal/data"
	"github.com/rarimo/points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateBalance(w http.ResponseWriter, r *http.Request) {
	did := r.Header.Get("X-User-DID")

	balance := getBalanceByDID(did, false, w, r)
	if balance != nil {
		ape.RenderErr(w, problems.Conflict())
		return
	}

	if err := BalancesQ(r).Insert(data.Balance{DID: did}); err != nil {
		Log(r).WithError(err).Error("Failed to create balance")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	// We can't return inserted balance in a single query, because we can't calculate
	// rank in transaction: RANK() is a window function allowed on a set of rows,
	// while with RETURNING we operate a single one.
	balance = getBalanceByDID(did, true, w, r)
	if balance == nil {
		return
	}

	err := EventsQ(r).Insert(prepareOpenEvents(balance.ID, EventTypes(r).List())...)
	if err != nil {
		Log(r).WithError(err).Error("Failed to add open events")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newBalanceModel(*balance))
}

func prepareOpenEvents(balanceID string, evTypes []resources.EventStaticMeta) []data.Event {
	events := make([]data.Event, len(evTypes))
	for i, evType := range evTypes {
		events[i] = data.Event{
			BalanceID: balanceID,
			Type:      evType.Name,
			Status:    data.EventOpen,
			PointsAmount: sql.NullInt32{
				Int32: evType.Reward,
				Valid: true,
			},
		}
	}
	return events
}
