package server

import "context"

type HandlerWrapper func(HandlerFunc) HandlerFunc

type HandlerFunc func(ctx context.Context, req Request, rsp interface{}) error
