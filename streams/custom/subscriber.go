package custom

import (
	pb "github.com/w-h-a/pkg/proto/streams"
	"github.com/w-h-a/pkg/streams"
)

type customSubscriber struct {
	options streams.SubscribeOptions
	channel chan streams.Event
}

func (s *customSubscriber) Options() streams.SubscribeOptions {
	return s.options
}

func (s *customSubscriber) Channel() chan streams.Event {
	return s.channel
}

func (s *customSubscriber) Close() {
	close(s.channel)
}

func (s *customSubscriber) Ack(ev streams.Event) interface{} {
	return &pb.AckRequest{Id: ev.Id, Success: true}
}

func (s *customSubscriber) Nack(ev streams.Event) interface{} {
	return &pb.AckRequest{Id: ev.Id, Success: false}
}

func (s *customSubscriber) SetAttempts(c int, ev streams.Event) {

}

func (s *customSubscriber) GetAttempts(ev streams.Event) (int, bool) {
	return 0, true
}

func (s *customSubscriber) String() string {
	return "custom"
}

func NewSubscriber(opts ...streams.SubscribeOption) streams.Subscriber {
	options := streams.NewSubscribeOptions(opts...)

	s := &customSubscriber{
		options: options,
		channel: make(chan streams.Event),
	}

	return s
}
