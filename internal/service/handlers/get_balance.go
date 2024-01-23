package handlers

import (
	"net/http"

	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
	did := r.Header.Get("X-User-DID") // TODO: get DID from auth

	balance := getBalanceByDID(did, true, w, r)
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
