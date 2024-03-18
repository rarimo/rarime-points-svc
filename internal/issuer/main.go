package issuer

import (
	"github.com/imroc/req/v3"
	"github.com/rarimo/rarime-points-svc/internal/config"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"strconv"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)

type Client struct {
	client *req.Client
	issuer string
	schema string
	typ    string
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		client: req.C().
			SetBaseURL(cfg.IssuerConfig().Host).
			SetCommonBasicAuth(cfg.IssuerConfig().Username, cfg.IssuerConfig().Password),
		schema: cfg.IssuerConfig().CredentialSchema,
		typ:    cfg.IssuerConfig().Type,
		issuer: cfg.IssuerConfig().Issuer,
	}
}

func (c *Client) IssueLevelClaim(did string, level int) (*string, error) {
	var result CreateCredentialResponse

	request := CreateCredentialRequest{
		CredentialSchema:  c.schema,
		Type:              c.typ,
		CredentialSubject: CredentialSubject{ID: did, Level: level},
		MtProof:           true,
		SignatureProof:    true,
	}

	response, err := c.client.R().
		SetBodyJsonMarshal(request).
		SetSuccessResult(&result).
		SetPathParam("identifier", c.issuer).
		Post("/{identifier}/claims")

	if err != nil {
		return nil, errors.Wrap(err, "failed to send post request")
	}

	if response.StatusCode >= 299 {
		return nil, errors.Wrap(ErrUnexpectedStatusCode, response.String())
	}

	return &result.ID, nil
}

func (c *Client) GetCredential(claimID string) (GetCredentialResponse, error) {
	var cred GetCredentialResponse

	response, err := c.client.R().
		SetSuccessResult(&cred).
		SetPathParam("id", claimID).
		Get("/credentials/{id}")
	if err != nil {
		return GetCredentialResponse{}, errors.Wrap(err, "failed to send post request")
	}

	if response.StatusCode >= 299 {
		return GetCredentialResponse{}, errors.Wrap(ErrUnexpectedStatusCode, response.String())
	}

	return cred, nil
}

func (c *Client) RevokeClaim(claimID string) error {
	credential, err := c.GetCredential(claimID)
	if err != nil {
		return errors.Wrap(err, "failed to get credential")
	}

	if !credential.Revoked {
		response, err := c.client.R().
			SetPathParam("nonce", strconv.FormatInt(credential.RevNonce, 10)).
			SetPathParam("identifier", c.issuer).
			Post("/{identifier}/claims/revoke/{nonce}")

		if err != nil {
			return errors.Wrap(err, "failed to send post request")
		}

		if response.StatusCode >= 299 {
			return errors.Wrap(ErrUnexpectedStatusCode, response.String())
		}
	}

	return nil
}

func (c *Client) GetClaimStatus(claimID string) (string, error) {
	var result GetClaimStateStatusResponse
	response, err := c.client.R().
		SetSuccessResult(&result).
		SetPathParam("identifier", c.issuer).
		SetPathParam("id", claimID).
		Get("/{identifier}/claims/{id}/status")
	if err != nil {
		return "", errors.Wrap(err, "failed to check claim's status")
	}
	defer response.Body.Close()

	return result.Status, nil
}
