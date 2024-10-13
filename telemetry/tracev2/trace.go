package tracev2

import (
	"context"
)

// TODO: status
type Trace interface {
	Options() TraceOptions
	// TODO: add cfg argument to check for whether tracing is enabled
	Start(ctx context.Context, name string) context.Context
	AddMetadata(ctx context.Context, md map[string]string)
	Finish(ctx context.Context)
	Read(opts ...ReadOption) ([]*SpanData, error)
	String() string
}
