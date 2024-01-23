package handlers

import (
	"net/http"

	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func Leaderboard(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewLeaderboard(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	leaders, err := BalancesQ(r).SelectLeaders(req.Limit)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance leaders")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newLeaderboardResponse(leaders))
}

func newLeaderboardResponse(balances []data.Balance) resources.BalanceListResponse {
	list := make([]resources.Balance, len(balances))
	for i, balance := range balances {
		list[i] = newBalanceModel(balance)
	}

	return resources.BalanceListResponse{Data: list}
}
