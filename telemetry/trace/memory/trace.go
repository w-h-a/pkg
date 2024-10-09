package memory

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/w-h-a/pkg/telemetry/buffer"
	"github.com/w-h-a/pkg/telemetry/buffer/memory"
	"github.com/w-h-a/pkg/telemetry/trace"
)

type memoryTrace struct {
	options trace.TraceOptions
	buffer  buffer.Buffer
}

func (t *memoryTrace) Options() trace.TraceOptions {
	return t.options
}

func (t *memoryTrace) Start(ctx context.Context, name string, md map[string]string) (context.Context, *trace.Span, error) {
	span := &trace.Span{
		Name:     name,
		Id:       uuid.New().String(),
		Trace:    uuid.New().String(),
		Started:  time.Now(),
		Metadata: md,
	}

	if ctx == nil {
		newCtx, err := trace.ContextWithIds(context.Background(), span.Trace, span.Id)
		return newCtx, span, err
	}

	traceId, traceFound, parentId, spanFound := trace.IdsFromContext(ctx)

	if !traceFound {
		newCtx, err := trace.ContextWithIds(ctx, span.Trace, span.Id)
		return newCtx, span, err
	} else {
		span.Trace = traceId
	}

	if spanFound {
		span.Parent = parentId
	}

	newCtx, err := trace.ContextWithIds(ctx, span.Trace, span.Id)
	return newCtx, span, err
}

func (t *memoryTrace) Finish(span *trace.Span) error {
	span.Duration = time.Since(span.Started)

	t.buffer.Put(span)

	return nil
}

func (t *memoryTrace) Read(opts ...trace.ReadOption) ([]*trace.Span, error) {
	options := trace.NewReadOptions(opts...)

	var entries []*buffer.Entry

	if options.Count > 0 {
		entries = t.buffer.Get(options.Count)
	} else {
		entries = t.buffer.Get(t.buffer.Options().Size)
	}

	spans := []*trace.Span{}

	for _, entry := range entries {
		span := entry.Value.(*trace.Span)

		if len(options.Trace) > 0 && options.Trace != span.Trace {
			continue
		}

		spans = append(spans, span)
	}

	return spans, nil
}

func (t *memoryTrace) String() string {
	return "memory"
}

func NewTrace(opts ...trace.TraceOption) trace.Trace {
	options := trace.NewTraceOptions(opts...)

	t := &memoryTrace{
		options: options,
	}

	if s, ok := GetSizeFromContext(options.Context); ok && s > 0 {
		t.buffer = memory.NewBuffer(buffer.BufferWithSize(s))
	} else {
		t.buffer = memory.NewBuffer()
	}

	return t
}
