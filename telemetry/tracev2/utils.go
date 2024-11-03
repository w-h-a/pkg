package tracev2

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/w-h-a/pkg/utils/metadatautils"
)

const (
	TraceParentKey = "traceparent"
)

func ContextWithTraceParent(ctx context.Context, traceparent string) (context.Context, error) {
	return metadatautils.MergeContext(ctx, map[string]string{
		TraceParentKey: traceparent,
	}, true), nil
}

func TraceIdFromContext(ctx context.Context) (traceId [16]byte, found bool) {
	traceparent, found := metadatautils.GetContext(ctx, TraceParentKey)
	if !found {
		return
	}

	parts := strings.Split(traceparent, "-")
	if len(parts) != 4 {
		return
	}

	decoded, err := hex.DecodeString(parts[1])
	if err != nil {
		return
	}

	copy(traceId[:], decoded)

	return
}

func SpanIdFromContext(ctx context.Context) (spanId [8]byte, found bool) {
	traceparent, found := metadatautils.GetContext(ctx, TraceParentKey)
	if !found {
		return
	}

	parts := strings.Split(traceparent, "-")
	if len(parts) != 4 {
		return
	}

	decoded, err := hex.DecodeString(parts[2])
	if err != nil {
		return
	}

	copy(spanId[:], decoded)

	return
}
