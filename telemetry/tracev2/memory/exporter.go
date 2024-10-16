package memory

import (
	"context"

	"github.com/w-h-a/pkg/telemetry/tracev2"
	"github.com/w-h-a/pkg/utils/memoryutils"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type memoryExporter struct {
	buffer *memoryutils.Buffer
}

func (e *memoryExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	spanData := []*tracev2.SpanData{}

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

		status := tracev2.Status{
			Code:        uint32(s.Status().Code),
			Description: s.Status().Description,
		}

		data := &tracev2.SpanData{
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
		e.buffer.Put(d)
	}

	return nil
}

// TODO?
func (e *memoryExporter) Shutdown(ctx context.Context) error {
	return nil
}

func NewExporter(buffer *memoryutils.Buffer) sdktrace.SpanExporter {
	return &memoryExporter{buffer}
}
