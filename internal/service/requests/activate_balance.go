package requests

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/resources"
)

func NewActivateBalance(r *http.Request) (req resources.CreateBalanceRequest, err error) {
	did := chi.URLParam(r, "did")
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":                     validation.Validate(req.Data.ID, validation.Required, validation.In(did)),
		"data/type":                   validation.Validate(req.Data.Type, validation.Required, validation.In(resources.UPDATE_BALANCE)),
		"data/attributes/referred_by": validation.Validate(req.Data.Attributes.ReferredBy, validation.Required),
	}

	return req, errs.Filter()
}
