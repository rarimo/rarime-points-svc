package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/resources"
)

func NewClaimEvent(r *http.Request) (req resources.Relation, err error) {
	id := chi.URLParam(r, "id")
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = fmt.Errorf("decode request body: %w", err)
		return
	}

	if req.Data == nil {
		err = validation.Errors{"data": validation.ErrRequired}
		return
	}

	return req, validation.Errors{
		"data/id":   validation.Validate(req.Data.ID, validation.Required, validation.In(id)),
		"data/type": validation.Validate(req.Data.Type, validation.Required, validation.In(resources.CLAIM_EVENT)),
	}.Filter()
}
