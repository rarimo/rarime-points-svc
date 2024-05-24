package handlers

import (
	"net/http"

	"github.com/rarimo/decentralized-auth-svc/pkg/auth"
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

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(req.Nullifier)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	var balance *data.Balance
	if req.Rank {
		balance, err = BalancesQ(r).GetWithRank(req.Nullifier)
	} else {
		balance, err = BalancesQ(r).FilterByNullifier(req.Nullifier).Get()
	}

	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by nullifier")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if balance == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var referrals []data.Referral
	if req.ReferralCodes {
		referrals, err = ReferralsQ(r).FilterByNullifier(req.Nullifier).Select()
		if err != nil {
			Log(r).WithError(err).Error("Failed to get referrals by nullifier")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	ape.Render(w, newBalanceResponse(*balance, referrals))
}

func newBalanceModel(balance data.Balance) resources.Balance {
	return resources.Balance{
		Key: resources.Key{
			ID:   balance.Nullifier,
			Type: resources.BALANCE,
		},
		Attributes: resources.BalanceAttributes{
			Amount:     balance.Amount,
			IsDisabled: !balance.ReferredBy.Valid,
			CreatedAt:  balance.CreatedAt,
			UpdatedAt:  balance.UpdatedAt,
			Rank:       balance.Rank,
			Level:      balance.Level,
		},
	}
}

func newBalanceResponse(balance data.Balance, referrals []data.Referral) resources.BalanceResponse {
	resp := resources.BalanceResponse{Data: newBalanceModel(balance)}
	if len(referrals) == 0 {
		return resp
	}

	activeCodes, consumedCodes := make([]string, 0, len(referrals)), make([]string, 0, len(referrals))
	resp.Data.Attributes.ActiveReferralCodes = &activeCodes
	resp.Data.Attributes.ConsumedReferralCodes = &consumedCodes

	for _, ref := range referrals {
		if ref.IsConsumed {
			consumedCodes = append(consumedCodes, ref.ID)
			continue
		}
		activeCodes = append(activeCodes, ref.ID)
	}

	return resp
}
