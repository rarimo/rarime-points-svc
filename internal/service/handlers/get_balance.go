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

func GetBalance(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetBalance(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	balance := getBalanceByDID(req.FilterDID, true, w, r)
	if balance == nil {
		return
	}

	ape.Render(w, newBalanceModel(*balance))
}

func newBalanceModel(balance data.Balance) resources.Balance {
	return resources.Balance{
		Key: resources.Key{
			ID:   balance.ID,
			Type: resources.BALANCE,
		},
		Attributes: resources.BalanceAttributes{
			Amount:    balance.Amount,
			UpdatedAt: balance.UpdatedAt,
			UserDid:   balance.DID,
			Rank:      balance.Rank,
		},
	}
}

func getBalanceByDID(did string, withRank bool, w http.ResponseWriter, r *http.Request) *data.Balance {
	if !auth.Authenticates(UserClaims(r), auth.UserGrant(did)) {
		ape.RenderErr(w, problems.Unauthorized())
		return nil
	}

	q := BalancesQ(r).FilterByUserDID(did)
	if withRank {
		q.WithRank()
	}

	balance, err := q.Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, problems.InternalError())
		return nil
	}

	if balance == nil {
		Log(r).Debugf("Balance not found for DID %s", did)
		ape.RenderErr(w, problems.NotFound())
		return nil
	}

	return balance
}
