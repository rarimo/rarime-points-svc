package connector

import (
	"net/http"
	"strconv"

	"github.com/google/jsonapi"
)

type FulfillEventRequest struct {
	Nullifier  string  `json:"nullifier"`
	EventType  string  `json:"event_type"`
	ExternalID *string `json:"external_id,omitempty"`
}

type FulfillVerifyProofEventRequest struct {
	Nullifier         string   `json:"nullifier"`
	ProofTypes        []string `json:"proof_types"`
	VerifierNullifier string   `json:"verifier_nullifier"`
}

type VerifyPassportRequest struct {
	Nullifier  string   `json:"nullifier"`
	Hash       string   `json:"hash"`
	SharedData []string `json:"shared_data"`
	IsUSA      bool     `json:"is_usa"`
}

// ErrorCode represents an error with a code indicating the unhappy flow that occurred
type ErrorCode string

const (
	CodeEventExpired     ErrorCode = "event_expired"     // event type is expired
	CodeEventDisabled    ErrorCode = "event_disabled"    // event type is disabled or not configured
	CodeEventNotFound    ErrorCode = "event_not_found"   // specific event not found for user
	CodeNullifierUnknown ErrorCode = "nullifier_unknown" // nullifier is unknown, while external_id was provided
	CodeInternalError    ErrorCode = "internal_error"    // other errors
)

func (c ErrorCode) JSONAPIError() *jsonapi.ErrorObject {
	status := http.StatusBadRequest
	if c == CodeInternalError {
		status = http.StatusInternalServerError
	}

	return &jsonapi.ErrorObject{
		Title:  http.StatusText(status),
		Status: strconv.Itoa(status),
		Code:   string(c),
	}
}

type Error struct {
	Code ErrorCode
	err  error
}

func (e *Error) Error() string {
	return e.err.Error()
}
