package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type EditReferralsRequest struct {
	Nullifier string  `json:"nullifier"`
	Count     *uint64 `json:"count"`
}

func NewEditReferrals(r *http.Request) (req EditReferralsRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	return req, validation.Errors{
		"nullifier": validation.Validate(req.Nullifier, validation.Required),
		"count":     validation.Validate(req.Count, validation.NotNil),
	}.Filter()
}
