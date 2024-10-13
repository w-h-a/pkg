package tracev2

import (
	"context"
	"errors"
)

var (
	ErrStart = errors.New("failed to start span")
)

// TODO: status
type Trace interface {
	Options() TraceOptions
	// TODO: add cfg argument to check for whether tracing is enabled
	Start(ctx context.Context, name string) (context.Context, Span, error)
	AddMetadata(span Span, md map[string]string)
	Finish(span Span)
	Read(opts ...ReadOption) ([]*SpanData, error)
	String() string
}
