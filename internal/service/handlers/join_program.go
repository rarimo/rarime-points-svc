package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func JoinProgram(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewJoinProgram(r)
	if err != nil {
		Log(r).WithError(err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	gotSig := r.Header.Get("Signature")
	wantSig := calculateCountrySignature(CountriesConfig(r).VerificationKey, req.Data.ID, req.Data.Attributes.Country)
	if gotSig != wantSig {
		Log(r).Warnf("Unauthorized access: HMAC signature mismatch: got %s, want %s", gotSig, wantSig)
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	balance, errs := getAndVerifyBalanceEligibility(r, req.Data.ID, nil)
	if len(errs) > 0 {
		ape.RenderErr(w, errs...)
		return
	}

	if balance.Country != nil {
		Log(r).Debugf("Balance %s already verified", balance.Nullifier)
		ape.RenderErr(w, problems.TooManyRequests())
		return
	}

	err = EventsQ(r).Transaction(func() error {
		return doPassportScanUpdates(r, *balance, req.Data.Attributes.Country)
	})
	if err != nil {
		Log(r).WithError(err).Error("Failed to execute transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	event, err := EventsQ(r).FilterByNullifier(balance.Nullifier).
		FilterByType(evtypes.TypePassportScan).
		FilterByStatus(data.EventClaimed).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get claimed event")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var res resources.PassportEventStateResponse
	res.Data.ID = req.Data.ID
	res.Data.Type = resources.PASSPORT_EVENT_STATE
	res.Data.Attributes.Claimed = event != nil

	ape.Render(w, res)
}

func calculateCountrySignature(key []byte, nullifier, country string) string {
	bNull, err := hex.DecodeString(nullifier)
	if err != nil {
		panic(fmt.Errorf("nullifier was not properly validated as hex: %w", err))
	}

	h := hmac.New(sha256.New, key)
	msg := append(bNull, []byte(country)...)
	h.Write(msg)

	return hex.EncodeToString(h.Sum(nil))
}
