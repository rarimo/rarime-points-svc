package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/resources"
)

func NewCreateBalance(r *http.Request) (req resources.CreateBalanceRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":   validation.Validate(req.Data.ID, validation.Required),
		"data/type": validation.Validate(req.Data.Type, validation.Required, validation.In(resources.CREATE_BALANCE)),
		"data/relationships/referral_link/data/id":   validation.Validate(req.Data.Relationships.ReferredBy.Data.ID, validation.Required),
		"data/relationships/referral_link/data/type": validation.Validate(req.Data.Relationships.ReferredBy.Data.Type, validation.Required, validation.In(resources.REFERRAL_CODE)),
	}

	return req, errs.Filter()
}
