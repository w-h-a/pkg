package memory

import (
	"github.com/w-h-a/pkg/telemetry/tracev2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type memorySpan struct {
	options tracev2.SpanOptions
	span    trace.Span
}

func (s *memorySpan) Options() tracev2.SpanOptions {
	return s.options
}

func (s *memorySpan) AddMetadata(md map[string]string) {
	attrs := []attribute.KeyValue{}

	for k, v := range md {
		attrs = append(attrs, attribute.String(k, v))
	}

	if len(attrs) > 0 {
		s.span.SetAttributes(attrs...)
	}
}

func (s *memorySpan) Finish() {
	s.span.End()
}

func (s *memorySpan) String() string {
	return "memory"
}

func NewSpan(opts ...tracev2.SpanOption) tracev2.Span {
	options := tracev2.NewSpanOptions(opts...)

	s := &memorySpan{
		options: options,
	}

	if span, ok := GetSpanFromContext(options.Context); ok {
		s.span = span
	}

	return s
}
