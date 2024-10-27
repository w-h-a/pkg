package otel

import (
	"context"
	"fmt"
	"sync"

	"github.com/w-h-a/pkg/telemetry/tracev2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type otelTrace struct {
	options tracev2.TraceOptions
	tracer  trace.Tracer
	spans   map[string]trace.Span
	mtx     sync.RWMutex
}

func (t *otelTrace) Options() tracev2.TraceOptions {
	return t.options
}

func (t *otelTrace) Start(ctx context.Context, name string) (context.Context, string) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	parentCtxCfg := trace.SpanContextConfig{}

	if spanId, ok := tracev2.SpanIdFromContext(ctx); ok {
		parentCtxCfg.SpanID = spanId
	}

	if traceId, ok := tracev2.TraceIdFromContext(ctx); ok {
		parentCtxCfg.TraceID = traceId
	}

	ctx, span := t.start(ctx, name, parentCtxCfg)

	key := span.SpanContext().SpanID().String()

	t.spans[key] = span

	newCtx, _ := tracev2.ContextWithTraceParent(ctx, fmt.Sprintf("00-%s-%s-00", span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String()))

	return newCtx, key
}

func (t *otelTrace) AddMetadata(span string, md map[string]string) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	if t.spans[span] == nil {
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

	t.spans[span].SetAttributes(attrs...)
}

func (t *otelTrace) UpdateStatus(span string, code uint32, description string) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	switch code {
	case 2:
		t.spans[span].SetStatus(codes.Ok, description)
	case 1:
		t.spans[span].SetStatus(codes.Error, description)
	default:
	}
}

func (t *otelTrace) Finish(span string) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.spans[span].End()

	delete(t.spans, span)
}

func (t *otelTrace) String() string {
	return "otel"
}

func (t *otelTrace) start(ctx context.Context, name string, parentCtxCfg trace.SpanContextConfig) (context.Context, trace.Span) {
	parentSpanCtx := trace.NewSpanContext(parentCtxCfg)

	ctx = trace.ContextWithRemoteSpanContext(ctx, parentSpanCtx)

	return t.tracer.Start(ctx, name)
}

func NewTrace(opts ...tracev2.TraceOption) tracev2.Trace {
	options := tracev2.NewTraceOptions(opts...)

	t := &otelTrace{
		options: options,
		tracer:  otel.Tracer(options.Name),
		spans:   map[string]trace.Span{},
		mtx:     sync.RWMutex{},
	}

	return t
}
