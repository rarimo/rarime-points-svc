package requests

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/internal/service/page"
	"gitlab.com/distributed_lab/urlval/v4"
)

type Leaderboard struct {
	page.OffsetParams
}

func NewLeaderboard(r *http.Request) (req Leaderboard, err error) {
	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		return req, validation.Errors{
			"query": fmt.Errorf("failed to decode query: %v", err),
		}
	}

	return req, req.Validate()
}
