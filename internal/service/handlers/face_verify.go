package handlers

import (
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

	if len(faceUserEvents) > 0 {
		log.Debugf("User has already verified face")
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

	newEvent := data.Event{
		Nullifier: nullifier,
		Type:      evtypes.TypeFaceParticipation,
		Status:    data.EventFulfilled,
		Meta:      data.Jsonb(fmt.Sprintf(`{"root_smt": "%s"}`, proof.PubSignals[RootSMT])),
	}

	if err = EventsQ(r).Insert(newEvent); err != nil {
		Log(r).WithError(err).Error("Failed to create face event")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
