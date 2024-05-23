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
		"nullifier":          req.Nullifier,
		"proof_types":        req.ProofTypes,
		"verifier_nullifier": req.VerifierNullifier,
	})

	owner, err := BalancesQ(r).FilterByNullifier(req.Nullifier).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get balance by nullifier")
		ape.RenderErr(w, api.CodeInternalError.JSONAPIError())
		return
	}

	// Normally should never happen
	if owner == nil {
		log.Error("Proof owner balance not exists")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	verifier, err := BalancesQ(r).FilterByNullifier(req.VerifierNullifier).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get verifier balance by nullifier")
		ape.RenderErr(w, api.CodeInternalError.JSONAPIError())
		return
	}

	// If the verifier does not have a balance, then create it
	if verifier == nil {
		events := EventTypes(r).PrepareEvents(req.VerifierNullifier, evtypes.FilterNotOpenable)
		typeExists := false
		for i, ev := range events {
			if eventTypeIsOneOfProofs(ev.Type, req.ProofTypes) {
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

		if err = createBalanceWithEvents(req.VerifierNullifier, "", events, r); err != nil {
			log.WithError(err).Error("Failed to create balance with events")
			ape.RenderErr(w, api.CodeInternalError.JSONAPIError())
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = EventsQ(r).Transaction(func() (err error) {
		passportValid := verifier.PassportHash.Valid && verifier.PassportExpires.Time.After(time.Now().UTC())
		if passportValid {
			log.Debugf("Verifier have valid passport.")
		}

		for _, proof := range req.ProofTypes {
			if err = verifyProofFulfill(r, req, req.VerifierNullifier, fmt.Sprintf("verify_proof_%s", proof)); err != nil {
				return
			}
			if !passportValid {
				continue
			}
			// The verifier must have a verified passport for the owner of the proof to receive points
			err = verifyProofFulfill(r, req, req.Nullifier, fmt.Sprintf("verified_proof_%s", proof))
			if err != nil {
				return
			}
		}

		return
	})
	if err != nil {
		log.WithError(err).Error("Failed to fulfill verify proof events")
		ape.RenderErr(w, api.CodeInternalError.JSONAPIError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func verifyProofFulfill(r *http.Request, req api.FulfillVerifyProofEventRequest, nullifier, evType string) (err error) {
	log := Log(r).WithFields(map[string]any{
		"nullifier":          req.Nullifier,
		"event_name":         evType,
		"verifier_nullifier": req.VerifierNullifier,
	})

	eventType := EventTypes(r).Get(evType, evtypes.FilterInactive)
	if eventType == nil {
		log.Debugf("Event %s inactive", evType)
		return nil
	}

	event, err := EventsQ(r).
		FilterByNullifier(nullifier).
		FilterByType(evType).
		FilterByStatus(data.EventOpen).
		Get()
	if err != nil {
		return fmt.Errorf("get event %s by nullifier: %w", evType, err)
	}

	if event == nil {
		log.Debugf("Event %s absent or already fulfilled for user", evType)
		return nil
	}

	_, err = EventsQ(r).FilterByID(event.ID).Update(data.EventFulfilled, nil, nil)
	if err != nil {
		return fmt.Errorf("update event %s by ID: %w", evType, err)
	}

	return nil
}

func eventTypeIsOneOfProofs(eventType string, proofs []string) bool {
	for _, proof := range proofs {
		if eventType == fmt.Sprintf("verify_proof_%s", proof) {
			return true
		}
	}

	return false
}
