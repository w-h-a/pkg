package grpcstreams

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/w-h-a/pkg/store"
	"github.com/w-h-a/pkg/streams"
	"github.com/w-h-a/pkg/telemetry/log"
)

type grpcStreams struct {
	options     streams.StreamsOptions
	store       store.Store
	subscribers map[string]streams.Subscriber
	mtx         sync.RWMutex
}

func (s *grpcStreams) Options() streams.StreamsOptions {
	return s.options
}

func (s *grpcStreams) Subscribe(id string, opts ...streams.SubscribeOption) error {
	sub := NewSubscriber(opts...)

	s.mtx.Lock()
	s.subscribers[id] = sub
	s.mtx.Unlock()

	return nil
}

func (s *grpcStreams) Unsubscribe(id string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.subscribers, id)

	return nil
}

func (s *grpcStreams) Consume(id string, opts ...streams.ConsumeOption) (<-chan streams.Event, error) {
	options := streams.NewConsumeOptions(opts...)

	s.mtx.RLock()
	sub, ok := s.subscribers[id]
	if !ok {
		s.mtx.RUnlock()
		return nil, streams.ErrSubscriberNotFound
	}
	s.mtx.RUnlock()

	if options.Offset.Unix() > 0 {
		go s.lookupPreviousEvents(sub, options.Offset)
	}

	return sub.Channel(), nil
}

func (s *grpcStreams) Produce(topic string, data interface{}, opts ...streams.ProduceOption) error {
	options := streams.NewProduceOptions(opts...)

	var payload []byte

	if p, ok := data.([]byte); ok {
		payload = p
	} else {
		p, err := json.Marshal(data)
		if err != nil {
			return streams.ErrEncodingData
		}
		payload = p
	}

	event := &streams.Event{
		Id:        uuid.New().String(),
		Topic:     topic,
		Payload:   payload,
		Timestamp: options.Timestamp,
		Metadata:  options.Metadata,
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		return streams.ErrEncodingEvent
	}

	key := fmt.Sprintf("%v:%v", event.Topic, event.Id)

	if err := s.store.Write(&store.Record{
		Key:   key,
		Value: bytes,
	}); err != nil {
		return fmt.Errorf("failed to write to event store: %v", err)
	}

	go s.handleEvent(event)

	return nil
}

func (s *grpcStreams) String() string {
	return "grpc"
}

func (s *grpcStreams) lookupPreviousEvents(sub streams.Subscriber, startTime time.Time) {
	recs, err := s.store.Read(sub.Options().Topic+":", store.ReadWithPrefix())
	if err != nil {
		log.Errorf("failed to find any previous events: %v", err)
		return
	}

	for _, rec := range recs {
		var event streams.Event

		if err := json.Unmarshal(rec.Value, &event); err != nil {
			continue
		}

		if event.Timestamp.Unix() < startTime.Unix() {
			continue
		}

		if err := SendEvent(sub, &event); err != nil {
			log.Errorf("failed to send previous event: %v", err)
			continue
		}
	}
}

func (s *grpcStreams) handleEvent(ev *streams.Event) {
	s.mtx.RLock()
	subs := s.subscribers
	s.mtx.RUnlock()

	groupedSubscribers := map[string]streams.Subscriber{}

	for _, sub := range subs {
		if len(sub.Options().Topic) == 0 || sub.Options().Topic == ev.Topic {
			groupedSubscribers[sub.Options().Group] = sub
		}
	}

	for _, sub := range groupedSubscribers {
		go func(sub streams.Subscriber) {
			if err := SendEvent(sub, ev); err != nil {
				log.Errorf("failed to handle event: %v", err)
			}
		}(sub)
	}
}

func NewStreams(opts ...streams.StreamsOption) streams.Streams {
	options := streams.NewStreamsOptions(opts...)

	g := &grpcStreams{
		options:     options,
		subscribers: map[string]streams.Subscriber{},
		mtx:         sync.RWMutex{},
	}

	s, ok := GetStoreFromContext(options.Context)
	if ok {
		g.store = s
	}

	return g
}
