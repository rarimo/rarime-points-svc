package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/pkg/connector"
)

func NewVerifyPassport(r *http.Request) (req connector.VerifyPassportRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	return req, validation.Errors{
		"nullifier":   validation.Validate(req.Nullifier, validation.Required),
		"hash":        validation.Validate(req.Hash, validation.Required),
		"shared_data": validation.Validate(req.SharedData, validation.Required, validation.Length(2, 0)),
	}.Filter()
}
