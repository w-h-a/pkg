package http

import (
	"context"
	"net/http"

	"github.com/w-h-a/pkg/serverv2"
)

type Middleware func(h http.Handler) http.Handler

type middlewareKey struct{}

func HttpServerWithMiddleware(ms ...Middleware) serverv2.ServerOption {
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
