package nats

import (
	"context"

	"github.com/w-h-a/pkg/broker"
)

type drainConnectionKey struct{}

func BrokerWithGracefulDisconnect() broker.BrokerOption {
	return func(o *broker.BrokerOptions) {
		o.Context = context.WithValue(o.Context, drainConnectionKey{}, true)
	}
}

func GetGracefulDisconnectFromContext(ctx context.Context) (bool, bool) {
	b, ok := ctx.Value(drainConnectionKey{}).(bool)
	return b, ok
}
