package handlers

import (
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
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
		ape.RenderErr(w, problems.BadRequest(errors.New("Balance inactive"))...)
		return
	}

	evType := EventTypes(r).Get(evtypes.TypePassportScan, evtypes.FilterInactive)
	if evType == nil {
		Log(r).Debug("Passport scan event absent, disabled, hasn't start yet or expired")
		w.WriteHeader(http.StatusNoContent)
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

	_, err = EventsQ(r).
		FilterByNullifier(nullifier).
		FilterByType(evtypes.TypePassportScan).
		Update(data.EventFulfilled, nil, nil)
	if err != nil {
		Log(r).WithError(err).Error("Failed to update passport scan event")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
