package requests

import (
	"net/http"

	"github.com/rarimo/rarime-points-svc/internal/service/page"
	"gitlab.com/distributed_lab/urlval/v4"
)

type Leaderboard struct {
	page.OffsetParams
	Count bool `url:"count"`
}

func NewLeaderboard(r *http.Request) (req Leaderboard, err error) {
	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		err = newDecodeError("query", err)
		return
	}

	return req, req.Validate()
}
