package requests

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/internal/service/page"
	"gitlab.com/distributed_lab/urlval/v4"
)

type ListWithdrawals struct {
	Nullifier string
	page.CursorParams
}

func NewListWithdrawals(r *http.Request) (req ListWithdrawals, err error) {
	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		err = newDecodeError("query", err)
		return
	}

	req.Nullifier = strings.ToLower(chi.URLParam(r, "nullifier"))
	return req, validation.Errors{
		"page":      req.Validate(),
		"nullifier": validation.Validate(req.Nullifier, validation.Required, validation.Match(nullifierRegexp)),
	}
}
