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
)

const privatePrefix = "/integrations/rarime-points-svc/v1/private"

type Client struct {
	conn    *conn.Connector
	ignored []ErrorCode
}

func NewClient(cli iface.Client) *Client {
	return &Client{conn: conn.NewConnector(cli)}
}

// IgnoreErrors creates a client copy which return nil error on the specified
// error codes. It simplifies error handling when you want to consider some cases
// a valid behaviour.
func (c *Client) IgnoreErrors(codes ...ErrorCode) *Client {
	return &Client{
		conn:    c.conn,
		ignored: codes,
	}
}

func (c *Client) FulfillEvent(ctx context.Context, req FulfillEventRequest) *Error {
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
	if c.isIgnored(code) {
		return nil
	}

	return &Error{
		Code: code,
		err:  baseErr,
	}
}

func (c *Client) VerifyPassport(ctx context.Context, req VerifyPassportRequest) error {
	u, _ := url.Parse(privatePrefix + "/balances")
	return c.conn.PatchJSON(u, req, ctx, nil)
}

func (c *Client) isIgnored(code ErrorCode) bool {
	for _, ig := range c.ignored {
		if ig == code {
			return true
		}
	}
	return false
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
