package http

import (
	"context"
	"net/http"

	"github.com/w-h-a/pkg/runner"
)

type handlerFuncsKey struct{}

func HttpProcessWithHandlerFuncs(route string, fun http.HandlerFunc) runner.ProcessOption {
	return func(o *runner.ProcessOptions) {
		funs := map[string]http.HandlerFunc{}

		if m, ok := GetHandlerFuncsFromContext(o.Context); ok && m != nil {
			m[route] = fun
			funs = m
		} else {
			funs[route] = fun
		}

		o.Context = context.WithValue(o.Context, handlerFuncsKey{}, funs)
	}
}

func GetHandlerFuncsFromContext(ctx context.Context) (map[string]http.HandlerFunc, bool) {
	funs, ok := ctx.Value(handlerFuncsKey{}).(map[string]http.HandlerFunc)
	return funs, ok
}
