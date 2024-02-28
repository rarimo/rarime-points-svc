package connector

import "time"

type FulfillEventRequest struct {
	UserDID    string  `json:"user_did"`
	EventType  string  `json:"event_type"`
	ExternalID *string `json:"external_id,omitempty"`
}

type VerifyPassportRequest struct {
	UserDID string    `json:"user_did"`
	Hash    string    `json:"hash"`
	Expiry  time.Time `json:"expiry"`
}
