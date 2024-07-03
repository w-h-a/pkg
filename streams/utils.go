package streams

import (
	"context"

	"github.com/w-h-a/pkg/utils/metadatautils"
)

const (
	subscriberKey = "streams-subscriber"
)

func ContextWithSubscriber(ctx context.Context, id string) (context.Context, error) {
	return metadatautils.SetContext(ctx, subscriberKey, id), nil
}

func SubscriberFromContext(ctx context.Context) (string, error) {
	str, ok := metadatautils.GetContext(ctx, subscriberKey)
	if !ok {
		return "", nil
	}

	return str, nil
}
