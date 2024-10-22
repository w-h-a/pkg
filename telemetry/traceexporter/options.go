package traceexporter

import (
	"context"

	"github.com/w-h-a/pkg/utils/memoryutils"
)

type ExporterOption func(o *ExporterOptions)

type ExporterOptions struct {
	Buffer   *memoryutils.Buffer
	Nodes    []string
	Protocol string
	Secure   bool
	Headers  map[string]string
	Context  context.Context
}

func ExporterWithBuffer(b *memoryutils.Buffer) ExporterOption {
	return func(o *ExporterOptions) {
		o.Buffer = b
	}
}

func ExporterWithNodes(addrs ...string) ExporterOption {
	return func(o *ExporterOptions) {
		o.Nodes = addrs
	}
}

func ExporterWithProtocol(p string) ExporterOption {
	return func(o *ExporterOptions) {
		o.Protocol = p
	}
}

func ExporterWithSecure() ExporterOption {
	return func(o *ExporterOptions) {
		o.Secure = true
	}
}

func ExporterWithHeaders(headers map[string]string) ExporterOption {
	return func(o *ExporterOptions) {
		o.Headers = headers
	}
}

func NewExporterOptions(opts ...ExporterOption) ExporterOptions {
	options := ExporterOptions{
		Nodes:   []string{},
		Headers: map[string]string{},
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
