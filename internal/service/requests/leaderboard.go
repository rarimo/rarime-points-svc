package requests

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/urlval/v4"
)

type Leaderboard struct {
	Limit int `page:"limit"`
}

func NewLeaderboard(r *http.Request) (req Leaderboard, err error) {
	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		return req, validation.Errors{
			"page[limit]": fmt.Errorf("failed to decode query: %v", err),
		}
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	err = validation.Errors{
		"page[limit]": validation.Validate(req.Limit, validation.Min(3), validation.Max(50)),
	}.Filter()

	return
}
