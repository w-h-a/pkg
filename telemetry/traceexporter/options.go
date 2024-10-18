package traceexporter

import "context"

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
