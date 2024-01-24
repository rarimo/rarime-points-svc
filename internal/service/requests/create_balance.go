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

	return req, validation.Errors{
		"data/id":                  validation.Validate(req.Data.ID, validation.Empty),
		"data/type":                validation.Validate(req.Data.Type, validation.Required, validation.In(resources.CREATE_BALANCE)),
		"data/attributes/user_did": validation.Validate(req.Data.Attributes.UserDid, validation.Required),
	}.Filter()
}
