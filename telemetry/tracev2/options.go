package tracev2

import "context"

type TraceOption func(o *TraceOptions)

type TraceOptions struct {
	Name    string
	Context context.Context
}

func TraceWithName(name string) TraceOption {
	return func(o *TraceOptions) {
		o.Name = name
	}
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

type SpanOption func(o *SpanOptions)

type SpanOptions struct {
	Name    string
	Context context.Context
}

func SpanWithName(name string) SpanOption {
	return func(o *SpanOptions) {
		o.Name = name
	}
}

func NewSpanOptions(opts ...SpanOption) SpanOptions {
	options := SpanOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type ReadOption func(o *ReadOptions)

type ReadOptions struct {
	Count   int
	Context context.Context
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
