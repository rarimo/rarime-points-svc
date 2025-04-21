package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/decentralized-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

const (
	RootSMT = iota
	FaceChallengedNullifier
)

func FaceVerify(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewFaceScanVerify(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	proof := req.Data.Attributes.Proof
	nullifier := UserClaims(r)[0].Nullifier

	log := Log(r).WithFields(map[string]any{
		"nullifier": nullifier,
		"proof":     proof,
	})

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(nullifier)) ||
		new(big.Int).SetBytes(hexutil.MustDecode(nullifier)).String() != proof.PubSignals[FaceChallengedNullifier] {
		log.Debug("failed to authenticate user")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	balance, err := BalancesQ(r).FilterByNullifier(nullifier).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get balance by nullifier")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if balance == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	evType := EventTypes(r).Get(evtypes.TypeFaceParticipation, evtypes.FilterInactive)
	if evType == nil {
		log.Infof("Event face participation type is inactive")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	faceUserEvents, err := EventsQ(r).FilterByNullifier(nullifier).FilterByType(evtypes.TypeFaceParticipation).Select()
	if err != nil {
		log.WithError(err).Error("Failed to get user face events")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// The event, if active, is opened for new users, so there must be at least one event.
	// If the no_auto_open parameter is set to true, the event will be created during the verification request only for these users.
	// If the no_auto_open parameter is false, then this event will be created for everyone and this logic will simply not work.
	if len(faceUserEvents) == 0 {
		log.Debugf("No face event found for nullifier=%s", nullifier)

		evNamesWithoutFaceEvent := EventTypes(r).Names(evtypes.FilterByNames(evtypes.TypeFaceParticipation))
		events := EventTypes(r).PrepareEvents(nullifier, evtypes.FilterNotOpenable, evtypes.FilterByNames(evNamesWithoutFaceEvent...))

		Log(r).Debugf("%d events will be added for nullifier=%s", len(events), nullifier)
		if err = EventsQ(r).Insert(events...); err != nil {
			Log(r).WithError(err).Error("Failed to create face event")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		faceUserEvents, err = EventsQ(r).FilterByNullifier(nullifier).FilterByType(evtypes.TypeFaceParticipation).Select()
		if err != nil {
			log.WithError(err).Error("Failed to get user face events")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	faceID := struct {
		RootSMT string `json:"root_smt"`
	}{}

	err = json.Unmarshal(faceUserEvents[0].Meta, &faceID)
	if err != nil {
		log.WithError(err).Errorf("Failed to parse event meta with eventID=%s", faceUserEvents[0].ID)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if faceID.RootSMT == proof.PubSignals[RootSMT] {
		log.Debugf("Face event already fulfilled")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	err = FaceVerifier(r).VerifyProof(proof)
	if err != nil {
		log.WithError(err).Debug("Failed to verify face participation proof")
		if errors.Is(err, config.ErrInvalidRoot) {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"proof": err,
			})...)
			return
		}

		ape.RenderErr(w, problems.InternalError())
		return
	}

	_, err = EventsQ(r).FilterByNullifier(nullifier).FilterByType(evtypes.TypeFaceParticipation).FilterByStatus(data.EventOpen).Update(
		data.EventFulfilled,
		json.RawMessage(fmt.Sprintf(`{"root_smt": "%s"}`, proof.PubSignals[RootSMT])),
		nil,
	)
	if err != nil {
		log.WithError(err).Error("Failed to insert poll event")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
