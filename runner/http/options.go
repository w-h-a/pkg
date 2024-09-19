package http

import (
	"context"
	"net/http"

	"github.com/w-h-a/pkg/runner"
)

type handlerFuncsKey struct{}

func ProcessWithHandlerFuncs(funs map[string]http.HandlerFunc) runner.ProcessOption {
	return func(o *runner.ProcessOptions) {
		o.Context = context.WithValue(o.Context, handlerFuncsKey{}, funs)
	}
}

func GetHandlerFuncsFromContext(ctx context.Context) (map[string]http.HandlerFunc, bool) {
	funs, ok := ctx.Value(handlerFuncsKey{}).(map[string]http.HandlerFunc)
	return funs, ok
}

type portKey struct{}

func ProcessWithPort(p int) runner.ProcessOption {
	return func(o *runner.ProcessOptions) {
		o.Context = context.WithValue(o.Context, portKey{}, p)
	}
}

func GetPortFromContext(ctx context.Context) (int, bool) {
	p, ok := ctx.Value(portKey{}).(int)
	return p, ok
}
