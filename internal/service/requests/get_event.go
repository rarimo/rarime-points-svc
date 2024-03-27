package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func NewGetEvent(r *http.Request) (id string, err error) {
	id = chi.URLParam(r, "id")
	return id, validation.Errors{"id": validation.Validate(id, validation.Required)}.
		Filter()
}
