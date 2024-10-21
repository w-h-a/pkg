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

func ExporterWithSecure(secure bool) ExporterOption {
	return func(o *ExporterOptions) {
		o.Secure = secure
	}
}

func NewExporterOptions(opts ...ExporterOption) ExporterOptions {
	options := ExporterOptions{
		Nodes:   []string{},
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
