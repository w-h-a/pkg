package memory

import (
	"context"

	"github.com/w-h-a/pkg/telemetry/traceexporter"
	"github.com/w-h-a/pkg/utils/memoryutils"
)

type bufferKey struct{}

func ExporterWithBuffer(b *memoryutils.Buffer) traceexporter.ExporterOption {
	return func(o *traceexporter.ExporterOptions) {
		o.Context = context.WithValue(o.Context, bufferKey{}, b)
	}
}

func GetBufferFromContext(ctx context.Context) (*memoryutils.Buffer, bool) {
	b, ok := ctx.Value(bufferKey{}).(*memoryutils.Buffer)
	return b, ok
}
