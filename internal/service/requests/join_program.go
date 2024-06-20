package requests

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rarimo/rarime-points-svc/resources"
)

func NewJoinProgram(r *http.Request) (req resources.JoinProgramRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	req.Data.ID = strings.ToLower(req.Data.ID)

	return req, validation.Errors{
		"data/id": validation.Validate(req.Data.ID,
			validation.Required,
			validation.In(strings.ToLower(chi.URLParam(r, "nullifier"))),
			validation.Match(nullifierRegexp)),
		"data/type": validation.Validate(req.Data.Type,
			validation.Required,
			validation.In(resources.JOIN_PROGRAM)),
		"data/attributes/country": validation.Validate(req.Data.Attributes.Country, validation.Required, is.CountryCode3),
	}.Filter()
}
