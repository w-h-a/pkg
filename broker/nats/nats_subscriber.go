package nats

import (
	client "github.com/nats-io/nats.go"
	"github.com/w-h-a/pkg/broker"
)

type natsSubscriber struct {
	options      broker.SubscribeOptions
	subscription *client.Subscription
}

func (s *natsSubscriber) Options() broker.SubscribeOptions {
	return s.options
}

func (s *natsSubscriber) Topic() string {
	return s.subscription.Subject
}

func (s *natsSubscriber) Unsubscribe() error {
	return s.subscription.Unsubscribe()
}

func (s *natsSubscriber) String() string {
	return "nats"
}
