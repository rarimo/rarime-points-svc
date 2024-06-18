package requests

import (
	"net/http"
	"strings"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/service/page"
	"gitlab.com/distributed_lab/urlval/v4"
)

type ListEvents struct {
	page.OffsetParams
	FilterNullifier     *string            `filter:"nullifier"`
	FilterStatus        []data.EventStatus `filter:"status"`
	FilterType          []string           `filter:"meta.static.name"`
	FilterHasExpiration *bool              `filter:"has_expiration"`
	FilterNotType       []string           `url:"filter[meta.static.name][not]"`
	Count               bool               `url:"count"`
}

func NewListEvents(r *http.Request) (req ListEvents, err error) {
	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		err = newDecodeError("query", err)
		return
	}
	if err = req.OffsetParams.Validate(); err != nil {
		return
	}

	if req.FilterNullifier != nil {
		*req.FilterNullifier = strings.ToLower(*req.FilterNullifier)
	}

	err = val.Errors{
		"filter[nullifier]":             val.Validate(req.FilterNullifier, val.Required, val.Match(nullifierRegexp)),
		"filter[status]":                val.Validate(req.FilterStatus, val.Each(val.In(data.EventOpen, data.EventFulfilled, data.EventClaimed))),
		"filter[meta.static.name][not]": val.Validate(req.FilterNotType, val.When(len(req.FilterType) > 0, val.Nil, val.Empty)),
	}.Filter()
	return
}
