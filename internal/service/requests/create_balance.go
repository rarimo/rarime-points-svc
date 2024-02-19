package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/resources"
)

func NewCreateBalance(r *http.Request) (req resources.CreateBalanceRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = fmt.Errorf("decode request body: %w", err)
		return
	}

	errs := validation.Errors{
		"data/id":   validation.Validate(req.Data.ID, validation.Required),
		"data/type": validation.Validate(req.Data.Type, validation.Required, validation.In(resources.CREATE_BALANCE)),
	}
	if attr := req.Data.Attributes; attr != nil {
		errs["data/attributes/referred_by"] = validation.Validate(attr.ReferredBy, validation.Required)
	}

	return req, errs.Filter()
}
