package rpc

import (
	"context"

	ctypes "github.com/cometbft/cometbft/rpc/core/types"
)

// Return false means the node finished syncing
func (c *Client) GetStatus(ctx context.Context) (*ctypes.ResultStatus, error) {
	resp, err := c.rpcClient.Status(ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
