package requests

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/resources"
)

var nullifierRegexp = regexp.MustCompile("^0x[0-9a-fA-F]{64}$")

func NewVerifyPassport(r *http.Request) (req resources.VerifyPassportRequest, err error) {
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
			validation.In(resources.VERIFY_PASSPORT)),
	}.Filter()
}
