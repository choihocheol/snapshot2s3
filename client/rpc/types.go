package rpc

import (
	cometbftHttp "github.com/cometbft/cometbft/rpc/client/http"
)

type Client struct {
	rpcClient *cometbftHttp.HTTP
	host string
}
