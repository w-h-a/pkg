package tracev2

import (
	"context"

	"github.com/w-h-a/pkg/utils/metadatautils"
)

const (
	TraceParentKey = "traceparent"
	SpanParentKey  = "spanparent"
)

func ContextWithTraceParent(ctx context.Context, traceparent [16]byte) (context.Context, error) {
	return metadatautils.MergeContext(ctx, map[string]string{
		TraceParentKey: string(traceparent[:]),
	}, true), nil
}

func TraceParentFromContext(ctx context.Context) (traceparent [16]byte, found bool) {
	traceId, found := metadatautils.GetContext(ctx, TraceParentKey)
	if !found {
		return
	}

	copy(traceparent[:], traceId)

	return
}

func ContextWithSpanParent(ctx context.Context, spanparent [8]byte) (context.Context, error) {
	return metadatautils.MergeContext(ctx, map[string]string{
		SpanParentKey: string(spanparent[:]),
	}, true), nil
}

func SpanParentFromContext(ctx context.Context) (spanparent [8]byte, found bool) {
	spanId, found := metadatautils.GetContext(ctx, SpanParentKey)
	if !found {
		return
	}

	copy(spanparent[:], spanId)

	return
}
