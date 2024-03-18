package connector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/jsonapi"
	conn "gitlab.com/distributed_lab/json-api-connector"
	"gitlab.com/distributed_lab/json-api-connector/cerrors"
	iface "gitlab.com/distributed_lab/json-api-connector/client"
	"gitlab.com/distributed_lab/logan/v3"
)

const privatePrefix = "/integrations/rarime-points-svc/v1/private"

type Client struct {
	disabled bool
	log      *logan.Entry
	conn     *conn.Connector
}

func NewClient(cli iface.Client) *Client {
	return &Client{conn: conn.NewConnector(cli), log: logan.New()}
}

func (c *Client) FulfillEvent(ctx context.Context, req FulfillEventRequest) *Error {
	if c.disabled {
		c.log.Info("Points connector disabled")
		return nil
	}

	u, _ := url.Parse(privatePrefix + "/events")

	err := c.conn.PatchJSON(u, req, ctx, nil)
	if err == nil {
		return nil
	}

	baseErr := err
	code, err := extractErrCode(err)
	if err != nil {
		return &Error{
			err: fmt.Errorf("failed to extract error code: %w; base error: %w", err, baseErr),
		}
	}

	return &Error{
		Code: code,
		err:  baseErr,
	}
}

func (c *Client) FulfillVerifyProofEvent(ctx context.Context, req FulfillVerifyProofEventRequest) *Error {
	if c.disabled {
		c.log.Info("Points connector disabled")
		return nil
	}

	u, _ := url.Parse(privatePrefix + "/proofs")

	err := c.conn.PatchJSON(u, req, ctx, nil)
	if err == nil {
		return nil
	}

	baseErr := err
	code, err := extractErrCode(err)
	if err != nil {
		return &Error{
			err: fmt.Errorf("failed to extract error code: %w; base error: %w", err, baseErr),
		}
	}

	return &Error{
		Code: code,
		err:  baseErr,
	}
}

func (c *Client) VerifyPassport(ctx context.Context, req VerifyPassportRequest) error {
	if c.disabled {
		c.log.Info("Points connector disabled")
		return nil
	}

	u, _ := url.Parse(privatePrefix + "/balances")
	return c.conn.PatchJSON(u, req, ctx, nil)
}

func extractErrCode(err error) (ErrorCode, error) {
	var apiErr cerrors.Error
	if !errors.As(err, &apiErr) {
		return "", errors.New("unknown error type")
	}

	var errs jsonapi.ErrorsPayload
	if errUn := json.Unmarshal(apiErr.Body(), &errs); errUn != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", errUn)
	}
	if len(errs.Errors) == 0 {
		return "", errors.New("empty errors payload")
	}

	return ErrorCode(errs.Errors[0].Code), nil
}
