package handlers

import (
	"net/http"
	"strconv"

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

	leaders, err := BalancesQ(r).Page(&req.CursorPageParams).Select()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance leaders")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var last string
	if len(leaders) > 0 {
		last = strconv.Itoa(leaders[len(leaders)-1].Amount)
	}

	resp := newLeaderboardResponse(leaders)
	resp.Links = req.GetCursorLinks(r, last)
	ape.Render(w, resp)
}

func newLeaderboardResponse(balances []data.Balance) resources.BalanceListResponse {
	list := make([]resources.Balance, len(balances))
	for i, balance := range balances {
		list[i] = newBalanceModel(balance)
	}

	return resources.BalanceListResponse{Data: list}
}
