package memory

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/w-h-a/pkg/broker"
)

type memoryBroker struct {
	options     broker.BrokerOptions
	subscribers map[string][]*memorySubscriber
	connected   bool
	addr        string
	mtx         sync.RWMutex
}

func (b *memoryBroker) Options() broker.BrokerOptions {
	return b.options
}

func (b *memoryBroker) Address() string {
	return b.addr
}

func (b *memoryBroker) Connect() error {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if b.connected {
		return nil
	}

	random := rand.Intn(20000)

	addr := fmt.Sprintf("127.0.0.1:%d", 10000+random)

	b.addr = addr

	b.connected = true

	return nil
}

func (b *memoryBroker) Disconnect() error {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if !b.connected {
		return nil
	}

	b.connected = false

	return nil
}

func (b *memoryBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	b.mtx.RLock()
	if !b.connected {
		b.mtx.RUnlock()
		return nil, broker.ErrNotConnected
	}
	b.mtx.RUnlock()

	options := broker.NewSubscribeOptions(opts...)

	subscriber := &memorySubscriber{
		options: options,
		id:      uuid.New().String(),
		topic:   topic,
		handler: handler,
		exit:    make(chan struct{}, 1),
	}

	b.mtx.Lock()
	b.subscribers[topic] = append(b.subscribers[topic], subscriber)
	b.mtx.Unlock()

	go func() {
		<-subscriber.exit

		b.mtx.Lock()

		newSubsForThisTopic := []*memorySubscriber{}
		for _, topicSubscriber := range b.subscribers[topic] {
			if topicSubscriber.id == subscriber.id {
				continue
			}
			newSubsForThisTopic = append(newSubsForThisTopic, topicSubscriber)
		}
		b.subscribers[topic] = newSubsForThisTopic

		b.mtx.Unlock()
	}()

	return subscriber, nil
}

func (b *memoryBroker) Publish(topic string, msg *broker.Message) error {
	b.mtx.RLock()

	if !b.connected {
		b.mtx.RUnlock()
		return broker.ErrNotConnected
	}

	subsOfThisTopic, ok := b.subscribers[topic]
	if !ok {
		b.mtx.RUnlock()
		return nil
	}

	b.mtx.RUnlock()

	publication := &memoryPublication{
		topic:   topic,
		message: msg,
	}

	for _, topicSubscriber := range subsOfThisTopic {
		if err := topicSubscriber.handler(publication); err != nil {
			return err
		}
	}

	return nil
}

func (b *memoryBroker) String() string {
	return "memory"
}

// just for random port
func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func NewBroker(opts ...broker.BrokerOption) broker.Broker {
	options := broker.NewBrokerOptions(opts...)

	b := &memoryBroker{
		options:     options,
		subscribers: make(map[string][]*memorySubscriber),
		mtx:         sync.RWMutex{},
	}

	return b
}
