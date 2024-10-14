package memory

import (
	"context"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/utils/memoryutils"
)

type bufferKey struct{}

func LogWithBuffer(b *memoryutils.Buffer) log.LogOption {
	return func(o *log.LogOptions) {
		o.Context = context.WithValue(o.Context, bufferKey{}, b)
	}
}

func GetBufferFromContext(ctx context.Context) (*memoryutils.Buffer, bool) {
	b, ok := ctx.Value(bufferKey{}).(*memoryutils.Buffer)
	return b, ok
}
