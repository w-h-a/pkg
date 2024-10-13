package memory

import (
	"context"
	"sync"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/telemetry/tracev2"
	"github.com/w-h-a/pkg/utils/memoryutils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type memoryTrace struct {
	options tracev2.TraceOptions
	tracer  trace.Tracer
	spans   map[string]trace.Span
	buffer  *memoryutils.Buffer
	mtx     sync.RWMutex
}

func (t *memoryTrace) Options() tracev2.TraceOptions {
	return t.options
}

func (t *memoryTrace) Start(ctx context.Context, name string) context.Context {
	// TODO: make sure we're enabled

	t.mtx.Lock()
	defer t.mtx.Unlock()

	parentCtxCfg := trace.SpanContextConfig{}

	if spanparent, ok := tracev2.SpanParentFromContext(ctx); ok {
		parentCtxCfg.SpanID = spanparent
	}

	if traceparent, ok := tracev2.TraceParentFromContext(ctx); ok {
		parentCtxCfg.TraceID = traceparent
	}

	ctx, span := t.start(ctx, name, parentCtxCfg)

	t.spans[span.SpanContext().SpanID().String()] = span

	newCtx, _ := tracev2.ContextWithSpan(ctx, span.SpanContext().SpanID())

	return newCtx
}

func (t *memoryTrace) AddMetadata(ctx context.Context, md map[string]string) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	span, found := tracev2.SpanFromContext(ctx)
	if !found {
		return
	}

	key := string(span[:])

	if t.spans[key] == nil {
		return
	}

	if len(md) == 0 {
		return
	}

	attrs := []attribute.KeyValue{}

	for k, v := range md {
		attrs = append(attrs, attribute.String(k, v))
	}

	if len(attrs) == 0 {
		return
	}

	t.spans[key].SetAttributes(attrs...)
}

func (t *memoryTrace) Finish(ctx context.Context) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	span, found := tracev2.SpanFromContext(ctx)
	if !found {
		return
	}

	key := string(span[:])

	t.spans[key].End()

	delete(t.spans, key)
}

func (t *memoryTrace) Read(opts ...tracev2.ReadOption) ([]*tracev2.SpanData, error) {
	options := tracev2.NewReadOptions(opts...)

	var entries []*memoryutils.Entry

	if options.Count > 0 {
		entries = t.buffer.Get(options.Count)
	} else {
		entries = t.buffer.Get(t.buffer.Size)
	}

	spans := []*tracev2.SpanData{}

	for _, entry := range entries {
		span := entry.Value.(*tracev2.SpanData)

		spans = append(spans, span)
	}

	return spans, nil
}

func (t *memoryTrace) String() string {
	return "memory"
}

func (t *memoryTrace) start(ctx context.Context, name string, parentCtxCfg trace.SpanContextConfig) (context.Context, trace.Span) {
	parentSpanCtx := trace.NewSpanContext(parentCtxCfg)

	ctx = trace.ContextWithRemoteSpanContext(ctx, parentSpanCtx)

	return t.tracer.Start(ctx, name)
}

func NewTrace(opts ...tracev2.TraceOption) tracev2.Trace {
	options := tracev2.NewTraceOptions(opts...)

	t := &memoryTrace{
		options: options,
		tracer:  otel.Tracer(options.Name),
		spans:   map[string]trace.Span{},
		mtx:     sync.RWMutex{},
	}

	if b, ok := GetBufferFromContext(options.Context); ok && b != nil {
		t.buffer = b
	} else {
		log.Fatalf("no buffer was given")
	}

	return t
}
