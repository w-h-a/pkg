package memory

import (
	"context"

	"github.com/w-h-a/pkg/telemetry/buffer"
	"github.com/w-h-a/pkg/telemetry/tracev2"
	"go.opentelemetry.io/otel/trace"
)

type bufferKey struct{}

func TraceWithBuffer(b buffer.Buffer) tracev2.TraceOption {
	return func(o *tracev2.TraceOptions) {
		o.Context = context.WithValue(o.Context, b, bufferKey{})
	}
}

func GetBufferFromContext(ctx context.Context) (buffer.Buffer, bool) {
	b, ok := ctx.Value(bufferKey{}).(buffer.Buffer)
	return b, ok
}

type spanKey struct{}

func SpanWithSpan(s trace.Span) tracev2.SpanOption {
	return func(o *tracev2.SpanOptions) {
		o.Context = context.WithValue(o.Context, s, spanKey{})
	}
}

func GetSpanFromContext(ctx context.Context) (trace.Span, bool) {
	s, ok := ctx.Value(spanKey{}).(trace.Span)
	return s, ok
}
