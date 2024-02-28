package connector

import (
	"context"
	"net/url"

	conn "gitlab.com/distributed_lab/json-api-connector"
	iface "gitlab.com/distributed_lab/json-api-connector/client"
)

const privatePrefix = "/integrations/rarime-points-svc/v1/private"

type Client struct {
	conn *conn.Connector
}

func NewClient(cli iface.Client) *Client {
	return &Client{conn: conn.NewConnector(cli)}
}

func (c *Client) FulfillEvent(ctx context.Context, req FulfillEventRequest) error {
	u, _ := url.Parse(privatePrefix + "/events")
	return c.conn.PatchJSON(u, req, ctx, nil)
}

func (c *Client) VerifyPassport(ctx context.Context, req VerifyPassportRequest) error {
	u, _ := url.Parse(privatePrefix + "/balances")
	return c.conn.PatchJSON(u, req, ctx, nil)
}
