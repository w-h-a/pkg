package grpc

import (
	"context"

	"github.com/w-h-a/pkg/serverv2"
)

type Middleware func(HandlerFunc) HandlerFunc

type HandlerFunc func(ctx context.Context, req interface{}, rsp interface{}) error

type middlewareKey struct{}

func GrpcServerWithMiddleware(ms ...Middleware) serverv2.ServerOption {
	return func(o *serverv2.ServerOptions) {
		middlewares := []Middleware{}

		if s, ok := GetMiddlewaresFromContext(o.Context); ok && s != nil {
			s = append(s, ms...)
			middlewares = s
		} else {
			middlewares = append(middlewares, ms...)
		}

		o.Context = context.WithValue(o.Context, middlewareKey{}, middlewares)
	}
}

func GetMiddlewaresFromContext(ctx context.Context) ([]Middleware, bool) {
	ms, ok := ctx.Value(middlewareKey{}).([]Middleware)
	return ms, ok
}
