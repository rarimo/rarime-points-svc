package requests

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/pkg/connector"
)

func NewVerifyPassport(r *http.Request) (req connector.VerifyPassportRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	return req, validation.Errors{
		"user_did": validation.Validate(req.UserDID, validation.Required),
		"hash":     validation.Validate(req.Hash, validation.Required),
		"expiry":   validation.Validate(req.Expiry, validation.Required, validation.By(isNotExpiredRule)),
	}.Filter()
}

func isNotExpiredRule(value interface{}) error {
	v, ok := value.(time.Time)
	if !ok {
		panic("value is not a time.Time") // invalid function usage
	}

	if v.Before(time.Now().UTC()) {
		return errors.New("expiry is in the past")
	}

	return nil
}
