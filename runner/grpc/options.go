package grpc

import (
	"context"

	"github.com/w-h-a/pkg/runner"
	"github.com/w-h-a/pkg/serverv2/grpc"
)

type handlerFuncsKey struct{}

func GrpcProcessWithHandlers(funs ...*grpc.Handler) runner.ProcessOption {
	return func(o *runner.ProcessOptions) {
		funs := []*grpc.Handler{}

		if m, ok := GetHandlersFromContext(o.Context); ok && m != nil {
			m = append(m, funs...)
			funs = m
		} else {
			funs = append(funs, funs...)
		}

		o.Context = context.WithValue(o.Context, handlerFuncsKey{}, funs)
	}
}

func GetHandlersFromContext(ctx context.Context) ([]*grpc.Handler, bool) {
	funs, ok := ctx.Value(handlerFuncsKey{}).([]*grpc.Handler)
	return funs, ok
}
