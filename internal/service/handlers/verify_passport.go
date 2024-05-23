package handlers

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	zkptypes "github.com/iden3/go-rapidsnark/types"
	"github.com/rarimo/decentralized-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	zk "github.com/rarimo/zkverifier-kit"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func VerifyPassport(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewVerifyPassport(r)
	if err != nil {
		Log(r).WithError(err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	nullifier := req.Data.ID

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(nullifier)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	balance, err := BalancesQ(r).FilterByNullifier(nullifier).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by nullifier")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if balance == nil {
		Log(r).Debug("Balance absent")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !balance.ReferredBy.Valid {
		Log(r).Debug("Balance inactive")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"referred_by": errors.New("balance inactive")})...)
		return
	}

	evType := EventTypes(r).Get(evtypes.TypePassportScan, evtypes.FilterInactive)
	if evType == nil {
		Log(r).Debug("Passport scan event absent, disabled, hasn't start yet or expired")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"passport_scan": errors.New("event disabled or absent")})...)
		return
	}

	event, err := EventsQ(r).FilterByNullifier(nullifier).
		FilterByType(evtypes.TypePassportScan).
		FilterByStatus(data.EventOpen).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get passport scan event")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if event == nil {
		Log(r).Debug("Event already fulfilled or absent for user")
		ape.RenderErr(w, problems.TooManyRequests())
		return
	}

	var proof zkptypes.ZKProof
	if err := json.Unmarshal(req.Data.Attributes.Proof, &proof); err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// MustDecode will never panic, because of the previous logic
	proof.PubSignals[zk.Nullifier] = new(big.Int).SetBytes(hexutil.MustDecode(nullifier)).String()
	if err := Verifier(r).VerifyProof(proof, zk.WithProofSelectorValue("23073")); err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var referralEvent = true
	evType = EventTypes(r).Get(evtypes.TypeReferralSpecific, evtypes.FilterInactive)
	if evType == nil {
		Log(r).Debug("Referral event type is disabled or expired, not accruing points to referrer")
		referralEvent = false
	}

	err = EventsQ(r).Transaction(func() (err error) {
		if referralEvent {
			// ReferredBy always valid because of the previous logic
			referral, err := ReferralsQ(r).Get(balance.ReferredBy.String)
			if err != nil {
				return errors.Wrap(err, "failed to get referral by ID")
			}

			err = EventsQ(r).Insert(data.Event{
				Nullifier: referral.Nullifier,
				Type:      evType.Name,
				Status:    data.EventFulfilled,
				Meta:      data.Jsonb(fmt.Sprintf(`{"nullifier": "%s"}`, nullifier)),
			})
			if err != nil {
				return errors.Wrap(err, "add event for referrer")
			}
		}

		_, err = EventsQ(r).
			FilterByID(event.ID).
			Update(data.EventFulfilled, nil, nil)
		if err != nil {
			return errors.Wrap(err, "failed to update passport scan event")
		}

		return nil
	})

	if err != nil {
		Log(r).WithError(err).Error("Failed to add referral event and update verify passport event")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
