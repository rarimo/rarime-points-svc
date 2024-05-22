package requests

import (
	"net/http"

	"github.com/go-chi/chi"
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

	req.Nullifier = chi.URLParam(r, "nullifier")
	return req, req.Validate()
}
