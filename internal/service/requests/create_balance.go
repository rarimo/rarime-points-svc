package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/resources"
)

func NewCreateBalance(r *http.Request) (req resources.Relation, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = fmt.Errorf("decode request body: %w", err)
		return
	}

	if req.Data == nil {
		err = validation.Errors{"data": validation.ErrRequired}
		return
	}

	return req, validation.Errors{
		"data/id":   validation.Validate(req.Data.ID, validation.Required),
		"data/type": validation.Validate(req.Data.Type, validation.Required, validation.In(resources.CREATE_BALANCE)),
	}.Filter()
}
