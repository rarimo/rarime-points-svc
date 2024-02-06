package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetBalance struct {
	DID string
}

func NewGetBalance(r *http.Request) (GetBalance, error) {
	did := chi.URLParam(r, "did")

	return GetBalance{did}, validation.Errors{
		"did": validation.Validate(did, validation.Required),
	}.Filter()
}
