package requests

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/urlval/v4"
)

type GetBalance struct {
	FilterDID string `filter:"did"`
}

func NewGetBalance(r *http.Request) (req GetBalance, err error) {
	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		return req, validation.Errors{
			"query": fmt.Errorf("failed to decode query: %w", err),
		}
	}

	err = validation.Errors{
		"filter[did]": validation.Validate(req.FilterDID, validation.Required),
	}.Filter()
	return
}
