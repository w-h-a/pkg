package memory

import (
	"context"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/telemetry/traceexporter"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type memoryExporter struct {
	options traceexporter.ExporterOptions
}

func (e *memoryExporter) Options() traceexporter.ExporterOptions {
	return e.options
}

func (e *memoryExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	spanData := []*traceexporter.SpanData{}

	for _, s := range spans {
		var parentSpanId trace.SpanID

		if s.Parent().IsValid() {
			parentSpanId = s.Parent().SpanID()
		}

		metadata := map[string]string{}

		for _, attr := range s.Attributes() {
			if attr.Value.Type() != attribute.STRING {
				continue
			}
			metadata[string(attr.Key)] = attr.Value.AsString()
		}

		status := traceexporter.Status{
			Code:        uint32(s.Status().Code),
			Description: s.Status().Description,
		}

		data := &traceexporter.SpanData{
			Name:     s.Name(),
			Id:       s.SpanContext().SpanID().String(),
			Parent:   parentSpanId.String(),
			Trace:    s.SpanContext().TraceID().String(),
			Started:  s.StartTime(),
			Ended:    s.EndTime(),
			Metadata: metadata,
			Status:   status,
		}

		spanData = append(spanData, data)
	}

	for _, d := range spanData {
		e.options.Buffer.Put(d)
	}

	return nil
}

// TODO: ?
func (e *memoryExporter) Shutdown(ctx context.Context) error {
	return nil
}

func (e *memoryExporter) String() string {
	return "memory"
}

func NewExporter(opts ...traceexporter.ExporterOption) traceexporter.TraceExporter {
	options := traceexporter.NewExporterOptions(opts...)

	if options.Buffer == nil {
		log.Fatalf("no buffer was given")
	}

	e := &memoryExporter{
		options: options,
	}

	return e
}
