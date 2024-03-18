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
		Post("/credentials")

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

	if credential.Revoked {
		return nil
	}

	response, err := c.client.R().
		SetPathParam("nonce", strconv.FormatInt(credential.RevNonce, 10)).
		Post("/credentials/revoke/{nonce}")

	if err != nil {
		return errors.Wrap(err, "failed to send post request")
	}

	if response.StatusCode >= 299 {
		return errors.Wrap(ErrUnexpectedStatusCode, response.String())
	}

	return nil
}
