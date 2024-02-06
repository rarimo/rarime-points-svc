package requests

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/internal/service/page"
	"gitlab.com/distributed_lab/urlval/v4"
)

type ListWithdrawals struct {
	DID string
	page.CursorParams
}

func NewListWithdrawals(r *http.Request) (req ListWithdrawals, err error) {
	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		return req, validation.Errors{
			"query": fmt.Errorf("failed to decode query: %w", err),
		}
	}

	req.DID = chi.URLParam(r, "did")
	return req, req.Validate()
}
