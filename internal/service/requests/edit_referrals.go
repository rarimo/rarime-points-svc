package requests

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/urlval/v4"
)

type EditReferralsRequest struct {
	DID   string  `url:"did"`
	Count *uint64 `url:"count"`
}

func NewEditReferrals(r *http.Request) (req EditReferralsRequest, err error) {
	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		err = newDecodeError("query", err)
		return
	}

	return req, validation.Errors{
		"did":   validation.Validate(req.DID, validation.Required),
		"count": validation.Validate(req.Count, validation.Required),
	}.Filter()
}
