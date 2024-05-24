package requests

import (
	"encoding/json"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/pkg/connector"
)

func NewFulfillEvent(r *http.Request) (req connector.FulfillEventRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	req.Nullifier = strings.ToLower(req.Nullifier)

	return req, validation.Errors{
		"nullifier":  validation.Validate(req.Nullifier, validation.Required, validation.Match(nullifierRegexp)),
		"event_type": validation.Validate(req.EventType, validation.Required),
	}.Filter()
}
