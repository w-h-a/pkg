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
