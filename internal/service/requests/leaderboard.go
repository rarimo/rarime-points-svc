package requests

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/internal/service/page"
	"gitlab.com/distributed_lab/urlval/v4"
)

const leaderboardDefaultLimit = 10

type Leaderboard struct {
	page.CursorParams
}

func NewLeaderboard(r *http.Request) (req Leaderboard, err error) {
	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		return req, validation.Errors{
			"page[limit]": fmt.Errorf("failed to decode query: %v", err),
		}
	}

	req.IsLeaderboard = true
	if req.Limit == 0 {
		req.Limit = leaderboardDefaultLimit
	}
	if err = req.Validate(); err != nil {
		return
	}

	err = validation.Errors{
		"page[limit]": validation.Validate(req.Limit, validation.Min(3), validation.Max(50)),
	}.Filter()

	return
}
