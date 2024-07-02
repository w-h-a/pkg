package grpcstreams

import (
	"sync"

	"github.com/w-h-a/pkg/streams"
)

type grpcSubscriber struct {
	options  streams.SubscribeOptions
	channel  chan streams.Event
	retryMap map[string]int
	mtx      sync.RWMutex
}

func (s *grpcSubscriber) Options() streams.SubscribeOptions {
	return s.options
}

func (s *grpcSubscriber) Channel() chan streams.Event {
	return s.channel
}

func (s *grpcSubscriber) Ack(ev streams.Event) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.retryMap, ev.Id)

	return nil
}

func (s *grpcSubscriber) Nack(ev streams.Event) error {
	return nil
}

func (s *grpcSubscriber) SetAttemptCount(c int, ev streams.Event) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.retryMap[ev.Id] = c
}

func (s *grpcSubscriber) GetAttemptCount(ev streams.Event) (int, bool) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	count, ok := s.retryMap[ev.Id]
	return count, ok
}

func (s *grpcSubscriber) String() string {
	return "grpc"
}

func NewSubscriber(opts ...streams.SubscribeOption) streams.Subscriber {
	options := streams.NewSubscribeOptions(opts...)

	s := &grpcSubscriber{
		options:  options,
		channel:  make(chan streams.Event),
		retryMap: map[string]int{},
		mtx:      sync.RWMutex{},
	}

	return s
}
