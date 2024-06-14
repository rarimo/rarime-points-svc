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
		referrals, err = ReferralsQ(r).
			FilterByNullifier(req.Nullifier).
			WithRewarding().
			Select()
		if err != nil {
			Log(r).WithError(err).Error("Failed to get referrals by nullifier with rewarding field")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	ape.Render(w, newBalanceResponse(*balance, referrals))
}

// newBalanceModel forms a balance response without referral fields, which must
// only be accessed with authorization.
func newBalanceModel(balance data.Balance) resources.Balance {
	return resources.Balance{
		Key: resources.Key{
			ID:   balance.Nullifier,
			Type: resources.BALANCE,
		},
		Attributes: resources.BalanceAttributes{
			Amount:    balance.Amount,
			CreatedAt: balance.CreatedAt,
			UpdatedAt: balance.UpdatedAt,
			Rank:      balance.Rank,
			Level:     balance.Level,
		},
	}
}

func newBalanceResponse(balance data.Balance, referrals []data.Referral) resources.BalanceResponse {
	resp := resources.BalanceResponse{Data: newBalanceModel(balance)}
	boolP := func(b bool) *bool { return &b }

	resp.Data.Attributes.IsDisabled = boolP(!balance.ReferredBy.Valid)
	resp.Data.Attributes.IsVerified = boolP(balance.Country != nil)

	if len(referrals) == 0 {
		return resp
	}

	var (
		active    = make([]string, 0, len(referrals))
		consumed  = make([]string, 0, len(referrals))
		rewarding = make([]string, 0, len(referrals))
	)
	resp.Data.Attributes.ActiveReferralCodes = &active
	resp.Data.Attributes.ConsumedReferralCodes = &consumed
	resp.Data.Attributes.RewardingReferralCodes = &rewarding

	for _, ref := range referrals {
		switch {
		case ref.UsageLeft > 0:
			active = append(active, ref.ID)
		case ref.IsRewarding:
			rewarding = append(rewarding, ref.ID)
		default:
			consumed = append(consumed, ref.ID)
		}
	}

	return resp
}
