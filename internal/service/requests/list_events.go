package requests

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/service/page"
	"gitlab.com/distributed_lab/urlval/v4"
)

type ListEvents struct {
	page.CursorParams
	FilterStatus []data.EventStatus `filter:"status"`
	Count        bool               `url:"count"`
}

func NewListEvents(r *http.Request) (req ListEvents, err error) {
	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		return req, validation.Errors{
			"query": fmt.Errorf("failed to decode query: %w", err),
		}
	}

	err = validation.Errors{
		"filter[status]": validation.Validate(req.FilterStatus, validation.Each(validation.In(data.EventOpen, data.EventFulfilled, data.EventClaimed))),
	}.Filter()
	return
}
