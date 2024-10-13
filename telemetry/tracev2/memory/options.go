package memory

import (
	"context"

	"github.com/w-h-a/pkg/telemetry/tracev2"
	"github.com/w-h-a/pkg/utils/memoryutils"
)

type bufferKey struct{}

func TraceWithBuffer(b *memoryutils.Buffer) tracev2.TraceOption {
	return func(o *tracev2.TraceOptions) {
		o.Context = context.WithValue(o.Context, bufferKey{}, b)
	}
}

func GetBufferFromContext(ctx context.Context) (*memoryutils.Buffer, bool) {
	b, ok := ctx.Value(bufferKey{}).(*memoryutils.Buffer)
	return b, ok
}
