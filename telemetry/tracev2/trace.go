package tracev2

import (
	"context"
)

// TODO: status
type Trace interface {
	Options() TraceOptions
	// TODO: add cfg argument to check for whether tracing is enabled
	Start(ctx context.Context, name string) (context.Context, string)
	AddMetadata(span string, md map[string]string)
	Finish(span string)
	Read(opts ...ReadOption) ([]*SpanData, error)
	String() string
}
