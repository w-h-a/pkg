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

type ExporterOption func(o *ExporterOptions)

type ExporterOptions struct {
	Context context.Context
}

func NewExporterOptions(opts ...ExporterOption) ExporterOptions {
	options := ExporterOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
