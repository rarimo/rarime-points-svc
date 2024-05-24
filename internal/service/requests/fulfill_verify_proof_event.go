package requests

import (
	"encoding/json"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/pkg/connector"
)

func NewFulfillVerifyProofEvent(r *http.Request) (req connector.FulfillVerifyProofEventRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	req.Nullifier = strings.ToLower(req.Nullifier)

	return req, validation.Errors{
		"nullifier":          validation.Validate(req.Nullifier, validation.Required, validation.Match(nullifierRegexp)),
		"proof_types":        validation.Validate(req.ProofTypes, validation.Required),
		"verifier_nullifier": validation.Validate(req.VerifierNullifier, validation.Required),
	}.Filter()
}
