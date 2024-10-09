package trace

import (
	"context"

	"github.com/w-h-a/pkg/utils/metadatautils"
)

const (
	traceIdKey = "trace-id"
	spanIdKey  = "span-id"
)

func ContextWithIds(ctx context.Context, traceId, parentId string) (context.Context, error) {
	return metadatautils.MergeContext(ctx, map[string]string{
		traceIdKey: traceId,
		spanIdKey:  parentId,
	}, true), nil
}

func IdsFromContext(ctx context.Context) (traceId string, foundTrace bool, parentId string, foundParent bool) {
	traceId, traceOk := metadatautils.GetContext(ctx, traceId)
	if !traceOk {
		return
	}

	parentId, spanOk := metadatautils.GetContext(ctx, spanIdKey)

	return traceId, traceOk, parentId, spanOk
}
