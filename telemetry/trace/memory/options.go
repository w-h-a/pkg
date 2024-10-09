package memory

import (
	"context"

	"github.com/w-h-a/pkg/telemetry/trace"
)

type sizeKey struct{}

func LogWithSize(s int) trace.TraceOption {
	return func(o *trace.TraceOptions) {
		o.Context = context.WithValue(o.Context, sizeKey{}, s)
	}
}

func GetSizeFromContext(ctx context.Context) (int, bool) {
	s, ok := ctx.Value(sizeKey{}).(int)
	return s, ok
}
