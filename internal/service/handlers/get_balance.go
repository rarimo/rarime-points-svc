package handlers

import (
	"net/http"
	"time"

	"github.com/rarimo/auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetBalance(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(req.DID)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	balance, err := BalancesQ(r).GetWithRank(req.DID)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if balance == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, resources.BalanceResponse{Data: newBalanceModel(*balance)})
}

func newBalanceModel(balance data.Balance) resources.Balance {
	return resources.Balance{
		Key: resources.Key{
			ID:   balance.DID,
			Type: resources.BALANCE,
		},
		Attributes: resources.BalanceAttributes{
			Amount:     balance.Amount,
			IsVerified: balance.PassportExpires.Time.After(time.Now().UTC()),
			CreatedAt:  balance.CreatedAt,
			UpdatedAt:  balance.UpdatedAt,
			Rank:       balance.Rank,
		},
	}
}
