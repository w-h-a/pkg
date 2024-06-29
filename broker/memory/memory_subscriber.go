package memory

import "github.com/w-h-a/pkg/broker"

type memorySubscriber struct {
	options broker.SubscribeOptions
	id      string
	topic   string
	handler broker.Handler
	exit    chan struct{}
}

func (s *memorySubscriber) Options() broker.SubscribeOptions {
	return s.options
}

func (s *memorySubscriber) Topic() string {
	return s.topic
}

func (s *memorySubscriber) Unsubscribe() error {
	s.exit <- struct{}{}
	return nil
}

func (s *memorySubscriber) String() string {
	return "memory"
}
