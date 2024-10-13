package tracev2

import (
	"context"
	"encoding/json"

	"github.com/w-h-a/pkg/utils/metadatautils"
)

const (
	spanDataKey = "trace-spanData"
)

func ContextWithSpanData(ctx context.Context, spanData *SpanData) (context.Context, error) {
	bytes, err := json.Marshal(spanData)
	if err != nil {
		return ctx, err
	}

	return metadatautils.SetContext(ctx, spanDataKey, string(bytes)), nil
}

func SpanDataFromContext(ctx context.Context) (*SpanData, error) {
	str, ok := metadatautils.GetContext(ctx, spanDataKey)
	if !ok {
		return nil, nil
	}

	s := &SpanData{}

	if err := json.Unmarshal([]byte(str), s); err != nil {
		return nil, err
	}

	return s, nil
}
