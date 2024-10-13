package memory

import (
	"context"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/telemetry/tracev2"
	"github.com/w-h-a/pkg/utils/memoryutils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type memoryTrace struct {
	options tracev2.TraceOptions
	tracer  trace.Tracer
	buffer  *memoryutils.Buffer
}

func (t *memoryTrace) Options() tracev2.TraceOptions {
	return t.options
}

func (t *memoryTrace) Start(ctx context.Context, name string) (context.Context, tracev2.Span, error) {
	// TODO: make sure we're enabled

	parentCtxCfg := trace.SpanContextConfig{}

	parentData, err := tracev2.SpanDataFromContext(ctx)
	if err != nil {
		newCtx, memorySpan := t.start(ctx, name, parentCtxCfg)
		return newCtx, memorySpan, err
	}

	parentCtxCfg.TraceID, err = trace.TraceIDFromHex(parentData.Trace)
	if err != nil {
		newCtx, memorySpan := t.start(ctx, name, parentCtxCfg)
		return newCtx, memorySpan, err
	}

	parentCtxCfg.SpanID, err = trace.SpanIDFromHex(parentData.Id)
	if err != nil {
		newCtx, memorySpan := t.start(ctx, name, parentCtxCfg)
		return newCtx, memorySpan, err
	}

	newCtx, memorySpan := t.start(ctx, name, parentCtxCfg)

	return newCtx, memorySpan, nil
}

func (t *memoryTrace) AddMetadata(span tracev2.Span, md map[string]string) {
	if span == nil {
		return
	}

	if len(md) > 0 {
		span.AddMetadata(md)
	}
}

func (t *memoryTrace) Finish(span tracev2.Span) {
	span.Finish()
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

func (t *memoryTrace) start(ctx context.Context, name string, parentCtxCfg trace.SpanContextConfig) (context.Context, tracev2.Span) {
	parentSpanCtx := trace.NewSpanContext(parentCtxCfg)

	ctx = trace.ContextWithRemoteSpanContext(ctx, parentSpanCtx)

	newCtx, span := t.tracer.Start(ctx, name)

	opts := []tracev2.SpanOption{
		tracev2.SpanWithName(name),
		SpanWithSpan(span),
	}

	memorySpan := NewSpan(opts...)

	return newCtx, memorySpan
}

func NewTrace(opts ...tracev2.TraceOption) tracev2.Trace {
	options := tracev2.NewTraceOptions(opts...)

	t := &memoryTrace{
		options: options,
		tracer:  otel.Tracer(options.Name),
	}

	b, ok := GetBufferFromContext(options.Context)

	log.Infof("MY CTX %+#v", options.Context)

	log.Infof("MY buffer %+#v", b)

	if ok && b != nil {
		t.buffer = b
	} else {
		log.Fatalf("no buffer was given")
	}

	return t
}
