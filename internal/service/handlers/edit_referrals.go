package handlers

import (
	"fmt"
	"net/http"

	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/service/referralid"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func EditReferrals(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewEditReferrals(r)
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

	if balance == nil {
		if *req.Count == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		events := prepareEventsWithRef(req.DID, "", r)
		if err = createBalanceWithEvents(req.DID, "", events, r); err != nil {
			Log(r).WithError(err).Error("Failed to create balance with events")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	var index uint64
	if balance != nil {
		index, err = ReferralsQ(r).FilterByUserDID(balance.DID).FilterByIsConsumed(false).Count()
		if err != nil {
			Log(r).WithError(err).Error("Failed to get referral count for user DID")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	added, err := adjustReferralsCount(index, req, r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to adjust referrals count")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, struct {
		Refs []string `json:"added_referrals"`
	}{added})
}

func prepareReferralsToAdd(did string, count, index uint64) []data.Referral {
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

func adjustReferralsCount(index uint64, req requests.EditReferralsRequest, r *http.Request) (refsAdded []string, err error) {
	switch {
	case *req.Count < index:
		toConsume := index - *req.Count
		if err = ReferralsQ(r).ConsumeFirst(req.DID, toConsume); err != nil {
			return nil, fmt.Errorf("consume referrals: %w", err)
		}
		Log(r).Infof("Consumed %d referrals for DID %s", toConsume, req.DID)

	case *req.Count > index:
		toAdd := *req.Count - index
		refsToAdd := prepareReferralsToAdd(req.DID, toAdd, index)
		if err = ReferralsQ(r).Insert(refsToAdd...); err != nil {
			return nil, fmt.Errorf("insert referrals: %w", err)
		}
		Log(r).Infof("Inserted %d referrals for DID %s", toAdd, req.DID)
		// while this is deterministic, the codes will be the same
		refsAdded = referralid.NewMany(req.DID, toAdd, index)

	default:
		Log(r).Infof("No referrals to add or consume for DID %s", req.DID)
	}

	return
}
