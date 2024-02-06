package handlers

import (
	"net/http"

	"github.com/rarimo/rarime-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func ListWithdrawals(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewListWithdrawals(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(req.DID)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	withdrawals, err := WithdrawalsQ(r).FilterByUserDID(req.DID).Page(&req.CursorPageParams).Select()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get withdrawal list")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var last string
	if len(withdrawals) > 0 {
		last = withdrawals[len(withdrawals)-1].ID
	}

	resp := newWithdrawalsResponse(withdrawals)
	resp.Links = req.CursorParams.GetLinks(r, last)
}

func newWithdrawalsResponse(withdrawals []data.Withdrawal) resources.WithdrawalListResponse {
	list := make([]resources.Withdrawal, len(withdrawals))
	for i, w := range withdrawals {
		list[i] = resources.Withdrawal{
			Key: resources.Key{
				ID:   w.ID,
				Type: resources.WITHDRAWAL,
			},
			Attributes: resources.WithdrawalAttributes{
				Amount:    w.Amount,
				Address:   w.Address,
				CreatedAt: w.CreatedAt,
			},
		}
	}
	return resources.WithdrawalListResponse{Data: list}
}
