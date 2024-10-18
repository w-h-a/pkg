package traceexporter

import (
	"context"

	"github.com/w-h-a/pkg/utils/memoryutils"
)

type ExporterOption func(o *ExporterOptions)

type ExporterOptions struct {
	Buffer  *memoryutils.Buffer
	Context context.Context
}

func ExporterWithBuffer(b *memoryutils.Buffer) ExporterOption {
	return func(o *ExporterOptions) {
		o.Buffer = b
	}
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
