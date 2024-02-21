package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/pkg/connector"
)

func NewFulfillEvent(r *http.Request) (req connector.FulfillEventRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	return req, validation.Errors{
		"user_did":   validation.Validate(req.UserDID, validation.Required),
		"event_type": validation.Validate(req.EventType, validation.Required),
	}
}
