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

	balance, err := BalancesQ(r).FilterByNullifier(req.Nullifier).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by nullifier")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if balance == nil {
		if *req.Count == 0 {
			Log(r).Debugf("Balance %s not found, skipping creation for count=0", req.Nullifier)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		events := prepareEventsWithRef(req.Nullifier, "", r)
		if err = createBalanceWithEvents(req.Nullifier, "", events, r); err != nil {
			Log(r).WithError(err).Error("Failed to create balance with events")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	if req.Genesis {
		count, err := ReferralsQ(r).FilterByNullifier(req.Nullifier).Count()
		if err != nil {
			Log(r).WithError(err).Error("Failed to get referrals count")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		referral := referralid.New(req.Nullifier, count)

		err = ReferralsQ(r).Insert(data.Referral{
			ID:        referral,
			Nullifier: req.Nullifier,
			UsageLeft: int32(*req.Count),
		})
		if err != nil {
			Log(r).WithError(err).Error("Failed to insert genesis referral")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		ape.Render(w, struct {
			Ref       string `json:"added_ref"`
			UsageLeft int    `json:"usage_left"`
		}{referral, int(*req.Count)})
		return
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

func prepareReferralsToAdd(nullifier string, count, index uint64) []data.Referral {
	refCodes := referralid.NewMany(nullifier, count, index)
	refs := make([]data.Referral, len(refCodes))

	for i, code := range refCodes {
		refs[i] = data.Referral{
			ID:        code,
			Nullifier: nullifier,
			UsageLeft: 1,
		}
	}

	return refs
}

func adjustReferralsCount(req requests.EditReferralsRequest, r *http.Request) (refsAdded []string, err error) {
	active, err := ReferralsQ(r).FilterByNullifier(req.Nullifier).FilterConsumed().Count()
	if err != nil {
		return nil, fmt.Errorf("count active referrals: %w", err)
	}

	if *req.Count == active {
		Log(r).Infof("No referrals to add or consume for nullifier %s", req.Nullifier)
		return
	}

	if *req.Count < active {
		toConsume := active - *req.Count
		if err = ReferralsQ(r).ConsumeFirst(req.Nullifier, toConsume); err != nil {
			return nil, fmt.Errorf("consume referrals: %w", err)
		}
		Log(r).Infof("Consumed %d referrals for nullifier %s", toConsume, req.Nullifier)
		return
	}

	index, err := ReferralsQ(r).FilterByNullifier(req.Nullifier).Count()
	if err != nil {
		return nil, fmt.Errorf("count all referrals: %w", err)
	}

	toAdd := *req.Count - active
	// balance must exist, according to preceding logic in EditReferrals
	err = ReferralsQ(r).Insert(prepareReferralsToAdd(req.Nullifier, toAdd, index)...)
	if err != nil {
		return nil, fmt.Errorf("insert referrals: %w", err)
	}
	Log(r).Infof("Inserted %d referrals for nullifier %s", toAdd, req.Nullifier)

	// while this is deterministic, the codes will be the same
	refsAdded = referralid.NewMany(req.Nullifier, toAdd, index)
	return
}
