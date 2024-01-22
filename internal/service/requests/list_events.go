package requests

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/points-svc/internal/data"
	"github.com/rarimo/points-svc/internal/service/page"
	"gitlab.com/distributed_lab/urlval/v4"
)

type ListEvents struct {
	page.CursorParams
	DID          string
	FilterStatus *data.EventStatus `filter:"status"`
}

func NewListEvents(r *http.Request) (req ListEvents, err error) {
	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		return req, validation.Errors{
			"query": fmt.Errorf("failed to decode query: %w", err),
		}
	}

	req.DID = r.Header.Get("X-User-DID")
	return
}
