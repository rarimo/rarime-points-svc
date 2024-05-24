package requests

import (
	"encoding/json"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/resources"
)

func NewCreateBalance(r *http.Request) (req resources.CreateBalanceRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	req.Data.ID = strings.ToLower(req.Data.ID)

	errs := validation.Errors{
		"data/id":                     validation.Validate(req.Data.ID, validation.Required, validation.Match(nullifierRegexp)),
		"data/type":                   validation.Validate(req.Data.Type, validation.Required, validation.In(resources.CREATE_BALANCE)),
		"data/attributes/referred_by": validation.Validate(req.Data.Attributes.ReferredBy, validation.Required),
	}

	return req, errs.Filter()
}
