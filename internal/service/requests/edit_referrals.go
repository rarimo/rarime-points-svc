package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type EditReferralsRequest struct {
	DID   string  `json:"did"`
	Count *uint64 `json:"count"`
}

func NewEditReferrals(r *http.Request) (req EditReferralsRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	return req, validation.Errors{
		"did":   validation.Validate(req.DID, validation.Required),
		"count": validation.Validate(req.Count, validation.NotNil),
	}.Filter()
}
