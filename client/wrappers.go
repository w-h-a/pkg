package client

import "context"

type ClientWrapper func(Client) Client

type CallWrapper func(CallFunc) CallFunc

type CallFunc func(ctx context.Context, address string, req Request, rsp interface{}, options CallOptions) (int, error)
