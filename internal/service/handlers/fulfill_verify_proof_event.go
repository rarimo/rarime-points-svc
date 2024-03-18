package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	api "github.com/rarimo/rarime-points-svc/pkg/connector"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func FulfillVerifyProofEvent(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewFulfillVerifyProofEvent(r)
	if err != nil {
		Log(r).WithError(err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	log := Log(r).WithFields(map[string]any{
		"user_did":     req.UserDID,
		"proof_type":   req.ProofType,
		"verifier_did": req.VerifierDID,
	})

	balance, err := BalancesQ(r).FilterByDID(req.UserDID).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, api.CodeInternalError.JSONAPIError())
		return
	}

	if balance == nil {
		log.Error("Balance not exists")
		ape.RenderErr(w, api.CodeInternalError.JSONAPIError())
		return
	}

	verifierBalance, err := BalancesQ(r).FilterByDID(req.VerifierDID).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get verifier balance by DID")
		ape.RenderErr(w, api.CodeInternalError.JSONAPIError())
		return
	}

	if verifierBalance == nil {
		log.Error("Verifier balance not exists")
		ape.RenderErr(w, api.CodeInternalError.JSONAPIError())
		return
	}

	if err = verifyProofFulfill(r, req, req.VerifierDID, fmt.Sprintf("verify_proof_%s", req.ProofType)); err != nil {
		log.WithError(err).Errorf("Failed to fulfill verify_proof_%s event for user", req.ProofType)
	}

	// The verifier must have a verified passport for the owner of the proof to receive points
	if verifierBalance.PassportHash.Valid && verifierBalance.PassportExpires.Time.Before(time.Now().UTC()) {
		if err = verifyProofFulfill(r, req, req.UserDID, fmt.Sprintf("verified_proof_%s", req.ProofType)); err != nil {
			log.WithError(err).Errorf("Failed to fulfill verified_proof_%s event for user", req.ProofType)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func verifyProofFulfill(r *http.Request, req api.FulfillVerifyProofEventRequest, did, eventName string) (err error) {
	eventType := EventTypes(r).Get(eventName, evtypes.FilterInactive)
	if eventType == nil {
		Log(r).WithFields(map[string]any{
			"user_did":     req.UserDID,
			"proof_type":   req.ProofType,
			"verifier_did": req.VerifierDID,
		}).Debugf("Event %s inactive", eventName)
		return nil
	}

	event, err := EventsQ(r).FilterByUserDID(did).
		FilterByType(eventName).
		FilterByStatus(data.EventOpen).Get()
	if err != nil {
		return fmt.Errorf("failed to get event %s by DID: %w", eventName, err)
	}

	if event == nil {
		Log(r).WithFields(map[string]any{
			"user_did":     req.UserDID,
			"proof_type":   req.ProofType,
			"verifier_did": req.VerifierDID,
		}).Debugf("Event %s absent or already fulfilled for user", eventName)
		return nil
	}

	_, err = EventsQ(r).FilterByID(event.ID).Update(data.EventFulfilled, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}
