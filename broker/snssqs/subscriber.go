package snssqs

import "github.com/w-h-a/pkg/broker"

type subscriber struct {
	options broker.SubscribeOptions
	id      string
	handler func([]byte) error
	exit    chan struct{}
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.options
}

func (s *subscriber) Id() string {
	return s.id
}

func (s *subscriber) Handler(b []byte) error {
	return s.handler(b)
}

func (s *subscriber) Unsubscribe() error {
	select {
	case <-s.exit:
		return nil
	default:
		close(s.exit)
		return nil
	}
}

func (s *subscriber) String() string {
	return "snssqs"
}
