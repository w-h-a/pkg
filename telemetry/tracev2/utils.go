package tracev2

import (
	"context"

	"github.com/w-h-a/pkg/utils/metadatautils"
)

const (
	traceParentKey = "traceparent"
)

func ContextWithTraceParent(ctx context.Context, traceparent [16]byte) (context.Context, error) {
	return metadatautils.MergeContext(ctx, map[string]string{
		traceParentKey: string(traceparent[:]),
	}, true), nil
}

func TraceParentFromContext(ctx context.Context) (traceparent [16]byte, found bool) {
	traceId, ok := metadatautils.GetContext(ctx, traceParentKey)
	if !ok {
		return
	}

	copy(traceparent[:], traceId)

	return
}
