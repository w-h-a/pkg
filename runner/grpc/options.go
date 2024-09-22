package grpc

import (
	"context"

	"github.com/w-h-a/pkg/runner"
	"github.com/w-h-a/pkg/serverv2/grpc"
)

type handlerFuncsKey struct{}

func GrpcProcessWithHandlers(hs ...*grpc.Handler) runner.ProcessOption {
	return func(o *runner.ProcessOptions) {
		handlers := []*grpc.Handler{}

		if s, ok := GetHandlersFromContext(o.Context); ok && s != nil {
			s = append(s, hs...)
			handlers = s
		} else {
			handlers = append(handlers, hs...)
		}

		o.Context = context.WithValue(o.Context, handlerFuncsKey{}, handlers)
	}
}

func GetHandlersFromContext(ctx context.Context) ([]*grpc.Handler, bool) {
	hs, ok := ctx.Value(handlerFuncsKey{}).([]*grpc.Handler)
	return hs, ok
}
