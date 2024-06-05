package server

import "context"

type ControllerWrapper func(ControllerFunc) ControllerFunc

type ControllerFunc func(ctx context.Context, req Request, rsp interface{}) error
