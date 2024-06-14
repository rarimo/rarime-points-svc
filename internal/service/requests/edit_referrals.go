package requests

import (
	"encoding/json"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type EditReferralsRequest struct {
	Nullifier string `json:"nullifier"`
	Count     uint64 `json:"count"`
}

func NewEditReferrals(r *http.Request) (req EditReferralsRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	req.Nullifier = strings.ToLower(req.Nullifier)

	return req, validation.Errors{
		"nullifier": validation.Validate(req.Nullifier, validation.Required, validation.Match(nullifierRegexp)),
		"count":     validation.Validate(req.Count, validation.Required),
	}.Filter()
}
