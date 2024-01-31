package handlers

import (
	"net/http"

	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateBalance(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewCreateBalance(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	did := req.Data.ID
	balance := getBalanceByDID(did, false, w, r)
	if balance != nil {
		ape.RenderErr(w, problems.Conflict())
		return
	}

	if err = BalancesQ(r).Insert(did); err != nil {
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

	err = EventsQ(r).Insert(EventTypes(r).PrepareOpenEvents(balance.DID)...)
	if err != nil {
		Log(r).WithError(err).Error("Failed to add open events")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newBalanceModel(*balance))
}
