package traceexporter

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type TraceExporter interface {
	sdktrace.SpanExporter
	Options() ExporterOptions
	String() string
}
