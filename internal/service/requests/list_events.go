package requests

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/service/page"
	"gitlab.com/distributed_lab/urlval/v4"
)

type ListEvents struct {
	page.OffsetParams
	FilterNullifier *string            `filter:"nullifier"`
	FilterStatus    []data.EventStatus `filter:"status"`
	FilterType      []string           `filter:"meta.static.name"`
	Count           bool               `url:"count"`
}

func NewListEvents(r *http.Request) (req ListEvents, err error) {
	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		err = newDecodeError("query", err)
		return
	}
	if err = req.OffsetParams.Validate(); err != nil {
		return
	}

	err = validation.Errors{
		"filter[nullifier]": validation.Validate(req.FilterNullifier, validation.Required),
		"filter[status]":    validation.Validate(req.FilterStatus, validation.Each(validation.In(data.EventOpen, data.EventFulfilled, data.EventClaimed))),
	}.Filter()
	return
}
