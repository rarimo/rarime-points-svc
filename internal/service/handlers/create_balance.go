package handlers

import (
	"fmt"
	"net/http"

	"github.com/rarimo/auth-svc/pkg/auth"
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

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(did)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	balance, err := getBalanceByDID(did, false, r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// Balance should not exist
	if balance != nil {
		ape.RenderErr(w, problems.Conflict())
		return
	}

	err = EventsQ(r).Transaction(func() error {
		if err = BalancesQ(r).Insert(did); err != nil {
			return fmt.Errorf("add balance: %w", err)
		}
		err = EventsQ(r).Insert(EventTypes(r).PrepareOpenEvents(balance.DID)...)
		if err != nil {
			return fmt.Errorf("add open events: %w", err)
		}
		return nil
	})
	if err != nil {
		Log(r).WithError(err).Error("Failed to add balance with open events")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// We can't return inserted balance in a single query, because we can't calculate
	// rank in transaction: RANK() is a window function allowed on a set of rows,
	// while with RETURNING we operate a single one.

	// Balance will exist cause of previous logic
	balance, err = getBalanceByDID(did, true, r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newBalanceModel(*balance))
}
