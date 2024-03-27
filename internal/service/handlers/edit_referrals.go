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
			Log(r).Debugf("Balance %s not found, skipping creation for count=0", req.DID)
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

	added, err := adjustReferralsCount(req, r)
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

func adjustReferralsCount(req requests.EditReferralsRequest, r *http.Request) (refsAdded []string, err error) {
	active, err := ReferralsQ(r).FilterByUserDID(req.DID).FilterByIsConsumed(false).Count()
	if err != nil {
		return nil, fmt.Errorf("count active referrals: %w", err)
	}

	if *req.Count == active {
		Log(r).Infof("No referrals to add or consume for DID %s", req.DID)
		return
	}

	if *req.Count < active {
		toConsume := active - *req.Count
		if err = ReferralsQ(r).ConsumeFirst(req.DID, toConsume); err != nil {
			return nil, fmt.Errorf("consume referrals: %w", err)
		}
		Log(r).Infof("Consumed %d referrals for DID %s", toConsume, req.DID)
		return
	}

	index, err := ReferralsQ(r).FilterByUserDID(req.DID).Count()
	if err != nil {
		return nil, fmt.Errorf("count all referrals: %w", err)
	}

	toAdd := *req.Count - active
	// balance must exist, according to preceding logic in EditReferrals
	err = ReferralsQ(r).Insert(prepareReferralsToAdd(req.DID, toAdd, index)...)
	if err != nil {
		return nil, fmt.Errorf("insert referrals: %w", err)
	}
	Log(r).Infof("Inserted %d referrals for DID %s", toAdd, req.DID)

	// while this is deterministic, the codes will be the same
	refsAdded = referralid.NewMany(req.DID, toAdd, index)
	return
}
