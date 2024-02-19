package connector

import (
	"context"
	"net/url"

	conn "gitlab.com/distributed_lab/json-api-connector"
	iface "gitlab.com/distributed_lab/json-api-connector/client"
)

const FulfillEventEndpoint = "/integrations/rarime-points-svc/v1/private/events"
const VerifyPassportEndpoint = "/integrations/rarime-points-svc/v1/private/balances"

type Client struct {
	conn *conn.Connector
}

func NewClient(cli iface.Client) *Client {
	return &Client{conn: conn.NewConnector(cli)}
}

func (c *Client) FulfillEvent(ctx context.Context, req FulfillEventRequest) error {
	u, _ := url.Parse(FulfillEventEndpoint)
	return c.conn.PostJSON(u, req, ctx, nil)
}

func (c *Client) VerifyPassport(ctx context.Context, req VerifyPassportRequest) error {
	u, _ := url.Parse(VerifyPassportEndpoint)
	return c.conn.PostJSON(u, req, ctx, nil)
}
