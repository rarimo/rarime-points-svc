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
		// WithoutExpiredStatus filters out referral codes that have an “expired” status.
		// The status “expired” is assigned to codes that have been used, but the party that used them did not complete the passport scanning procedure by the set time.
		// In this case, usage_left is set to -1, and new codes are generated to replace the expired ones.
		// This allows you to keep a history of all codes used by the user.
		// It can be used only after applying the WithStatus filter, since the statuses are defined in it.
		referrals, err = ReferralsQ(r).
			FilterByNullifier(req.Nullifier).
			WithStatus().
			WithoutExpiredStatus().
			Select()
		if err != nil {
			Log(r).WithError(err).Error("Failed to get referrals by nullifier with rewarding field")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	ape.Render(w, newBalanceResponse(*balance, referrals))
}

// NewBalanceModel forms a balance response without referral fields, which must
// only be accessed with authorization.
func NewBalanceModel(balance data.Balance) resources.Balance {
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
	resp := resources.BalanceResponse{Data: NewBalanceModel(balance)}
	boolP := func(b bool) *bool { return &b }

	resp.Data.Attributes.IsDisabled = boolP(!balance.ReferredBy.Valid)
	resp.Data.Attributes.IsVerified = boolP(balance.Country != nil)

	if len(referrals) == 0 {
		return resp
	}

	res := make([]resources.ReferralCode, len(referrals))
	for i, r := range referrals {
		res[i] = resources.ReferralCode{
			Id:     r.ID,
			Status: r.Status,
		}
	}

	resp.Data.Attributes.ReferralCodes = &res
	return resp
}
