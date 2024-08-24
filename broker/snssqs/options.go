package snssqs

import (
	"context"

	"github.com/w-h-a/pkg/broker"
)

type snsClientKey struct{}
type sqsClientKey struct{}

func SnsSqsWithSnsClient(c SnsClient) broker.BrokerOption {
	return func(o *broker.BrokerOptions) {
		o.Context = context.WithValue(o.Context, snsClientKey{}, c)
	}
}

func GetSnsClientFromContext(ctx context.Context) (SnsClient, bool) {
	c, ok := ctx.Value(snsClientKey{}).(SnsClient)
	return c, ok
}

func SnsSqsWithSqsClient(c SqsClient) broker.BrokerOption {
	return func(o *broker.BrokerOptions) {
		o.Context = context.WithValue(o.Context, sqsClientKey{}, c)
	}
}

func GetSqsClientFromContext(ctx context.Context) (SqsClient, bool) {
	c, ok := ctx.Value(sqsClientKey{}).(SqsClient)
	return c, ok
}

type visibilityTimeoutKey struct{}
type waitTimeSecondsKey struct{}

func SqsWithVisibilityTimeout(t int32) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		o.Context = context.WithValue(o.Context, visibilityTimeoutKey{}, t)
	}
}

func GetVisibilityTimeoutFromContext(ctx context.Context) (int32, bool) {
	t, ok := ctx.Value(visibilityTimeoutKey{}).(int32)
	return t, ok
}

func SqsWithWaitTimeSeconds(t int32) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		o.Context = context.WithValue(o.Context, waitTimeSecondsKey{}, t)
	}
}

func GetWaitTimeSecondsFromContext(ctx context.Context) (int32, bool) {
	t, ok := ctx.Value(waitTimeSecondsKey{}).(int32)
	return t, ok
}
