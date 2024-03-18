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

	owner, err := BalancesQ(r).FilterByDID(req.UserDID).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, api.CodeInternalError.JSONAPIError())
		return
	}

	// Normally should never happen
	if owner == nil {
		log.Error("Proof owner balance not exists")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	verifier, err := BalancesQ(r).FilterByDID(req.VerifierDID).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get verifier balance by DID")
		ape.RenderErr(w, api.CodeInternalError.JSONAPIError())
		return
	}

	// If the verifier does not have a balance, then create it
	if verifier == nil {
		events := EventTypes(r).PrepareEvents(req.VerifierDID, evtypes.FilterNotOpenable)
		typeExists := false
		for i, ev := range events {
			if ev.Type == fmt.Sprintf("verify_proof_%s", req.ProofType) {
				events[i].Status = data.EventFulfilled
				typeExists = true
				break
			}
		}

		if !typeExists {
			log.Debug("Event type is not openable")
			ape.RenderErr(w, api.CodeEventNotFound.JSONAPIError())
			return
		}

		if err = createBalanceWithEvents(req.VerifierDID, "", events, r); err != nil {
			log.WithError(err).Error("Failed to create balance with events")
			ape.RenderErr(w, api.CodeInternalError.JSONAPIError())
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = EventsQ(r).Transaction(func() (err error) {
		if err = verifyProofFulfill(r, req, req.VerifierDID, fmt.Sprintf("verify_proof_%s", req.ProofType)); err != nil {
			return
		}

		// The verifier must have a verified passport for the owner of the proof to receive points
		if verifier.PassportHash.Valid && verifier.PassportExpires.Time.Before(time.Now().UTC()) {
			log.Debugf("Verifier have valid passport. Fulfill event for proof owner")
			return verifyProofFulfill(r, req, req.UserDID, fmt.Sprintf("verified_proof_%s", req.ProofType))
		}

		log.Debugf("Verifier haven't valid passport")
		return
	})
	if err != nil {
		log.WithError(err).Error("Failed to fulfill verify proof events")
		ape.RenderErr(w, api.CodeInternalError.JSONAPIError())
		return
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
