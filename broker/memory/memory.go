package memory

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/w-h-a/pkg/broker"
)

type memory struct {
	options     broker.BrokerOptions
	subscribers map[string][]broker.Subscriber
	mtx         sync.RWMutex
}

func (b *memory) Options() broker.BrokerOptions {
	return b.options
}

func (b *memory) Publish(data interface{}, options broker.PublishOptions) error {
	b.mtx.RLock()
	subsOfThisTopic, ok := b.subscribers[options.Topic]
	if !ok {
		b.mtx.RUnlock()
		return nil
	}
	b.mtx.RUnlock()

	var bs []byte

	if p, ok := data.([]byte); ok {
		bs = p
	} else {
		p, err := json.Marshal(data)
		if err != nil {
			return err
		}
		bs = p
	}

	for _, sub := range subsOfThisTopic {
		if err := sub.Handler(bs); err != nil {
			return err
		}
	}

	return nil
}

func (b *memory) Subscribe(callback func([]byte) error, options broker.SubscribeOptions) broker.Subscriber {
	b.mtx.Lock()

	sub := &subscriber{
		options: options,
		id:      uuid.New().String(),
		handler: callback,
		exit:    make(chan struct{}, 1),
	}

	b.subscribers[options.Group] = append(b.subscribers[options.Group], sub)

	b.mtx.Unlock()

	go func() {
		<-sub.exit

		b.mtx.Lock()

		newSubsForThisGroup := []broker.Subscriber{}

		for _, subscriber := range b.subscribers[options.Group] {
			if subscriber.Id() == sub.id {
				continue
			}
			newSubsForThisGroup = append(newSubsForThisGroup, subscriber)
		}

		b.subscribers[options.Group] = newSubsForThisGroup

		b.mtx.Unlock()
	}()

	return sub
}

func (b *memory) String() string {
	return "memory"
}

func NewBroker(opts ...broker.BrokerOption) broker.Broker {
	options := broker.NewBrokerOptions(opts...)

	b := &memory{
		options:     options,
		subscribers: map[string][]broker.Subscriber{},
		mtx:         sync.RWMutex{},
	}

	return b
}
