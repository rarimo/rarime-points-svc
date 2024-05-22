package requests

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/resources"
)

func NewWithdraw(r *http.Request) (req resources.WithdrawRequest, err error) {
	nullifier := chi.URLParam(r, "nullifier")

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	return req, validation.Errors{
		"data/id":                 validation.Validate(req.Data.ID, validation.Required, validation.In(nullifier)),
		"data/type":               validation.Validate(req.Data.Type, validation.Required, validation.In(resources.WITHDRAW)),
		"data/attributes/amount":  validation.Validate(req.Data.Attributes.Amount, validation.Required, validation.Min(1)),
		"data/attributes/address": validation.Validate(req.Data.Attributes.Address, validation.Required),
	}.Filter()
}
