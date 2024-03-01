package handlers

import (
	"net/http"

	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/service/referralid"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func AddReferrals(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewAddReferrals(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	balance, err := BalancesQ(r).FilterByDID(req.DID).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var index uint
	if balance == nil {
		events := prepareEventsWithRef(req.DID, "", r)
		if err = createBalanceWithEvents(req.DID, "", events, r); err != nil {
			Log(r).WithError(err).Error("Failed to create balance with events")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	} else {
		index, err = ReferralsQ(r).FilterByUserDID(balance.DID).Count()
		if err != nil {
			Log(r).WithError(err).Error("Failed to get referral count for user DID")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	refCodes := referralid.NewMany(req.DID, req.Count, index)
	refs := make([]data.Referral, len(refCodes))
	for i, code := range refCodes {
		refs[i] = data.Referral{
			ID:      code,
			UserDID: req.DID,
		}
	}

	err = ReferralsQ(r).Insert(prepareReferralsToAdd(req.DID, req.Count, index)...)
	if err != nil {
		Log(r).WithError(err).Error("Failed to insert referrals")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// TODO: return balance WITHOUT rank and with referrals included
	w.WriteHeader(http.StatusNoContent)
}

func prepareReferralsToAdd(did string, count, index uint) []data.Referral {
	refCodes := referralid.NewMany(did, count, index)
	refs := make([]data.Referral, len(refCodes))

	for i, code := range refCodes {
		refs[i] = data.Referral{
			ID:      code,
			UserDID: did,
		}
	}

	return refs
}
