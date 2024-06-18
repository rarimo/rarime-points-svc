package requests

import (
	"net/http"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"gitlab.com/distributed_lab/urlval/v4"
)

type ListExpiredEvents struct {
	FilterName    []string `filter:"name"`
	FilterFlag    []string `filter:"flag"`
	FilterNotName []string `url:"filter[name][not]"`
}

func NewListEventTypes(r *http.Request) (req ListExpiredEvents, err error) {
	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		err = newDecodeError("query", err)
		return
	}

	err = val.Errors{
		"filter[flag]": val.Validate(req.FilterFlag, val.Each(val.In(
			evtypes.FlagActive,
			evtypes.FlagNotStarted,
			evtypes.FlagExpired,
			evtypes.FlagDisabled,
		))),
		"filter[name][not]": val.Validate(req.FilterNotName, val.When(len(req.FilterName) > 0, val.Nil, val.Empty)),
	}.Filter()

	return
}
