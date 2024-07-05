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

func Leaderboard(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewLeaderboard(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	leaders, err := BalancesQ(r).FilterDisabled().Page(&req.OffsetPageParams).SelectWithRank()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance leaders")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	leadersCount, err := BalancesQ(r).FilterDisabled().Count()
	if err != nil {
		Log(r).WithError(err).Error("Failed to leaders count")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	resp := newLeaderboardResponse(leaders, leadersCount)
	resp.Links = req.GetLinks(r)
	ape.Render(w, resp)
}

func newLeaderboardResponse(balances []data.Balance, totalLeaders int64) resources.BalanceListResponse {
	list := make([]resources.Balance, len(balances))
	for i, balance := range balances {
		list[i] = newBalanceModel(balance)
	}

	// must never panic
	serMeta, _ := json.Marshal(struct {
		TotalCount int64 `json:"total_count"`
	}{
		TotalCount: totalLeaders,
	})

	return resources.BalanceListResponse{Data: list, Meta: serMeta}
}
