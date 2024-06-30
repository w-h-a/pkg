package nats

import (
	"context"

	"github.com/w-h-a/pkg/broker"
)

type gracefulDisconnectKey struct{}

func BrokerWithGracefulDisconnect() broker.BrokerOption {
	return func(o *broker.BrokerOptions) {
		o.Context = context.WithValue(o.Context, gracefulDisconnectKey{}, true)
	}
}

func GetGracefulDisconnectFromContext(ctx context.Context) (bool, bool) {
	b, ok := ctx.Value(gracefulDisconnectKey{}).(bool)
	return b, ok
}
