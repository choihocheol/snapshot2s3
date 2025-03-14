package rpc

import (
	"context"

	cometbftHttp "github.com/cometbft/cometbft/rpc/client/http"
)

func New(host string) (*Client, error) {
	c := &Client{
		host: host,
	}
	client, err := cometbftHttp.New(c.host)
	if err != nil {
		return nil, err
	}
	c.rpcClient = client

	return c, nil
}

func (c *Client) Connect(ctx context.Context) error {
	// For websocket connection
	err := c.rpcClient.Start()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Terminate(_ context.Context) error {
	// For websocket connection
	err := c.rpcClient.Stop()
	if err != nil {
		return err
	}

	return nil
}
