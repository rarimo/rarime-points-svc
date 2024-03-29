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

	var balance *data.Balance
	if req.Rank {
		balance, err = BalancesQ(r).GetWithRank(req.DID)
	} else {
		balance, err = BalancesQ(r).FilterByDID(req.DID).Get()
	}

	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if balance == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var referrals []data.Referral
	if req.ReferralCodes {
		referrals, err = ReferralsQ(r).FilterByUserDID(req.DID).FilterByIsConsumed(false).Select()
		if err != nil {
			Log(r).WithError(err).Error("Failed to get referrals by DID")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	ape.Render(w, newBalanceResponse(*balance, referrals))
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
			IsDisabled: !balance.ReferredBy.Valid,
			CreatedAt:  balance.CreatedAt,
			UpdatedAt:  balance.UpdatedAt,
			Rank:       balance.Rank,
		},
	}
}

func newBalanceResponse(balance data.Balance, referrals []data.Referral) resources.BalanceResponse {
	balanceResponse := resources.BalanceResponse{Data: newBalanceModel(balance)}
	if len(referrals) == 0 {
		return balanceResponse
	}

	referralCodes := make([]string, len(referrals))
	balanceResponse.Data.Attributes.ReferralCodes = &referralCodes
	for i, referral := range referrals {
		referralCodes[i] = referral.ID
	}
	return balanceResponse
}
