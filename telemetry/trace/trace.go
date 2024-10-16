package trace

import (
	"context"
	"errors"
)

var (
	ErrStart = errors.New("failed to start span")
)

var tracer Trace

type Trace interface {
	Options() TraceOptions
	Start(ctx context.Context, name string, md map[string]string) (context.Context, *Span, error)
	Finish(span *Span) error
	Read(opts ...ReadOption) ([]*Span, error)
	String() string
}

func SetTracer(t Trace) {
	tracer = t
}

func GetTracer() Trace {
	return tracer
}
