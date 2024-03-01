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
		index, err = ReferralsQ(r).FilterByUserDID(balance.DID).Count()
		if err != nil {
			Log(r).WithError(err).Error("Failed to get referral count for user DID")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	if err = adjustReferralsCount(index, req, r); err != nil {
		Log(r).WithError(err).Error("Failed to adjust referrals count")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// TODO: return balance WITHOUT rank and with referrals included, or just referrals, also above
	w.WriteHeader(http.StatusNoContent)
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

func adjustReferralsCount(index uint64, req requests.EditReferralsRequest, r *http.Request) error {
	switch {
	case *req.Count < index:
		toConsume := index - *req.Count
		if err := ReferralsQ(r).ConsumeFirst(req.DID, toConsume); err != nil {
			return fmt.Errorf("consume referrals: %w", err)
		}
		Log(r).Infof("Consumed %d referrals for DID %s", toConsume, req.DID)

	case *req.Count > index:
		toAdd := *req.Count - index
		err := ReferralsQ(r).Insert(prepareReferralsToAdd(req.DID, toAdd, index)...)
		if err != nil {
			return fmt.Errorf("insert referrals: %w", err)
		}
		Log(r).Infof("Inserted %d referrals for DID %s", toAdd, req.DID)

	default:
		Log(r).Infof("No referrals to add or consume for DID %s", req.DID)
	}

	return nil
}
