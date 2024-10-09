package trace

import "context"

type TraceOption func(o *TraceOptions)

type TraceOptions struct {
	Context context.Context
}

func NewTraceOptions(opts ...TraceOption) TraceOptions {
	options := TraceOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type ReadOption func(o *ReadOptions)

type ReadOptions struct {
	Trace   string
	Count   int
	Context context.Context
}

func ReadWithTrace(t string) ReadOption {
	return func(o *ReadOptions) {
		o.Trace = t
	}
}

func ReadWithCount(c int) ReadOption {
	return func(o *ReadOptions) {
		o.Count = c
	}
}

func NewReadOptions(opts ...ReadOption) ReadOptions {
	options := ReadOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
