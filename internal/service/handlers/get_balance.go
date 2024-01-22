package handlers

import (
	"net/http"

	"github.com/rarimo/points-svc/internal/data"
	"github.com/rarimo/points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
	did := r.Header.Get("X-User-DID") // TODO: get DID from auth
	balance := getBalanceByDID(did, w, r)
	if balance == nil {
		return
	}

	ape.Render(w, resources.BalanceResponse{
		Data: resources.Balance{
			Key: resources.Key{
				ID:   balance.ID,
				Type: resources.BALANCE,
			},
			Attributes: resources.BalanceAttributes{
				Amount:    balance.Amount,
				UpdatedAt: balance.UpdatedAt,
				UserDid:   balance.DID,
			},
		},
	})
}

func getBalanceByDID(did string, w http.ResponseWriter, r *http.Request) *data.Balance {
	balance, err := BalancesQ(r).FilterByUserDID(did).Get()
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
