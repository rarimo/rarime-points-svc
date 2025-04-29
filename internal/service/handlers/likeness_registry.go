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
	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

const (
	RootSMT = iota
	RootChallengedNullifier
)

func LiklessRegistry(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewLikenessRegistryVerify(r)
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
		new(big.Int).SetBytes(hexutil.MustDecode(nullifier)).String() != proof.PubSignals[RootChallengedNullifier] {
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

	evType := EventTypes(r).Get(evtypes.TypeLikenessRegistry, evtypes.FilterInactive)
	if evType == nil {
		log.Infof("Event likeness registry type is inactive")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	userEventsRootInclusion, err := EventsQ(r).FilterByNullifier(nullifier).FilterByType(evtypes.TypeLikenessRegistry).Select()
	if err != nil {
		log.WithError(err).Error("Failed to get user likeness registry events")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if len(userEventsRootInclusion) > 0 {
		log.Debugf("User has already verified likeness in registry")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	err = RootInclusionVerifier(r).VerifyProof(proof)
	if err != nil {
		log.WithError(err).Debug("Failed to verify likeness proof")
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
		Type:      evtypes.TypeLikenessRegistry,
		Status:    data.EventFulfilled,
		Meta:      data.Jsonb(fmt.Sprintf(`{"root_smt": "%s"}`, proof.PubSignals[RootSMT])),
	}

	if err = EventsQ(r).Insert(newEvent); err != nil {
		Log(r).WithError(err).Error("Failed to create likeness event")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	event, err := EventsQ(r).FilterByNullifier(balance.Nullifier).
		FilterByType(evtypes.TypeLikenessRegistry).
		FilterByStatus(data.EventFulfilled).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get fulfilled event")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, resources.LikenessRegistryEventState{
		Key: resources.Key{
			ID:   event.ID,
			Type: resources.LIKENESS_REGISTRY_EVENT_STATE,
		},
		Attributes: resources.LikenessRegistryEventStateAttributes{
			Fulfilled: event.Status == data.EventFulfilled,
		},
	})
}
