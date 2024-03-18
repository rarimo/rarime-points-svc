package issuer

import (
	"encoding/json"
	"time"
)

type CreateCredentialRequest struct {
	CredentialSchema  string            `json:"credentialSchema"`
	CredentialSubject CredentialSubject `json:"credentialSubject"`
	Expiration        *time.Time        `json:"expiration,omitempty"`
	MtProof           bool              `json:"mtProof,omitempty"`
	SignatureProof    bool              `json:"signatureProof,omitempty"`
	Type              string            `json:"type"`
}

type CredentialSubject struct {
	ID    string `json:"id"`
	Level int    `json:"level"`
}

type CreateCredentialResponse struct {
	ID string `json:"id"`
}

type GetClaimStateStatusResponse struct {
	Status string `json:"status"`
}

type GetCredentialResponse struct {
	Id                    string          `json:"id"`
	ProofTypes            []string        `json:"proofTypes"`
	CreatedAt             time.Time       `json:"createdAt"`
	ExpiresAt             time.Time       `json:"expiresAt"`
	Expired               bool            `json:"expired"`
	SchemaHash            string          `json:"schemaHash"`
	SchemaType            string          `json:"schemaType"`
	SchemaUrl             string          `json:"schemaUrl"`
	Revoked               bool            `json:"revoked"`
	RevNonce              int64           `json:"revNonce"`
	CredentialSubject     json.RawMessage `json:"credentialSubject"`
	UserID                string          `json:"userID"`
	SchemaTypeDescription string          `json:"schemaTypeDescription"`
}
