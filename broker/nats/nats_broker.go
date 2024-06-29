package nats

import (
	"strings"
	"sync"

	client "github.com/nats-io/nats.go"
	"github.com/w-h-a/pkg/broker"
	"github.com/w-h-a/pkg/telemetry/log"
)

type natsBroker struct {
	options     broker.BrokerOptions
	natsOptions client.Options
	addrs       []string
	connected   bool
	connection  *client.Conn
	drain       bool
	close       chan error
	mtx         sync.RWMutex
}

func (b *natsBroker) Options() broker.BrokerOptions {
	return b.options
}

func (b *natsBroker) Address() string {
	if b.connection != nil && b.connection.IsConnected() {
		return b.connection.ConnectedUrl()
	}

	if len(b.addrs) > 0 {
		return b.addrs[0]
	}

	return ""
}

func (b *natsBroker) Connect() error {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if b.connected {
		return nil
	}

	status := client.CLOSED

	if b.connection != nil {
		status = b.connection.Status()
	}

	switch status {
	case client.CONNECTED, client.RECONNECTING, client.CONNECTING:
		b.connected = true

		return nil
	default: // DISCONNECTED, CLOSED, DRAINING
		natsOptions := b.natsOptions

		natsOptions.Servers = b.addrs

		natsOptions.TLSConfig = b.options.TLSConfig

		if b.options.TLSConfig != nil || b.options.Secure {
			natsOptions.Secure = true
		}

		connection, err := natsOptions.Connect()
		if err != nil {
			return err
		}

		b.connection = connection

		b.connected = true

		return nil
	}
}

func (b *natsBroker) Disconnect() error {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if b.drain {
		b.connection.Drain()
		b.close <- nil
	}

	b.connection.Close()

	b.connected = false

	return nil
}

func (b *natsBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	b.mtx.RLock()
	if b.connection == nil {
		b.mtx.RUnlock()
		return nil, broker.ErrNotConnected
	}
	b.mtx.RUnlock()

	options := broker.NewSubscribeOptions(opts...)

	// wrap the handler to make it work with nats
	fn := func(msg *client.Msg) {
		message := broker.Message{}

		if err := b.options.Marshaler.Unmarshal(msg.Data, &message); err != nil {
			message.Body = msg.Data
			log.Errorf("failed to unmarshal message: %v", err)
			return
		}

		publication := &natsPublication{
			topic:   msg.Subject,
			message: &message,
		}

		if err := handler(publication); err != nil {
			log.Errorf("failed to handle publication: %v", err)
			return
		}
	}

	var subscription *client.Subscription

	var err error

	b.mtx.RLock()

	if len(options.QueueName) > 0 {
		subscription, err = b.connection.QueueSubscribe(topic, options.QueueName, fn)
	} else {
		subscription, err = b.connection.Subscribe(topic, fn)
	}

	b.mtx.RUnlock()

	if err != nil {
		return nil, err
	}

	subscriber := &natsSubscriber{
		options:      options,
		subscription: subscription,
	}

	return subscriber, nil
}

func (b *natsBroker) Publish(topic string, msg *broker.Message) error {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	if b.connection == nil {
		return broker.ErrNotConnected
	}

	bytes, err := b.options.Marshaler.Marshal(msg)
	if err != nil {
		return err
	}

	return b.connection.Publish(topic, bytes)
}

func (b *natsBroker) String() string {
	return "nats"
}

func (b *natsBroker) configure() error {
	b.natsOptions = client.GetDefaultOptions()

	if len(b.options.Nodes) == 0 {
		b.options.Nodes = b.natsOptions.Servers
	}

	if !b.options.Secure {
		b.options.Secure = b.natsOptions.Secure
	}

	if b.options.TLSConfig == nil {
		b.options.TLSConfig = b.natsOptions.TLSConfig
	}

	b.addrs = b.setAddrs()

	if graceful, ok := GetGracefulDisconnectFromContext(b.options.Context); ok && graceful {
		b.drain = true
		b.close = make(chan error)
		b.natsOptions.ClosedCB = b.onClose
		b.natsOptions.AsyncErrorCB = b.onAsyncError
		b.natsOptions.DisconnectedErrCB = b.onDisconnecedError
	}

	return nil
}

func (b *natsBroker) setAddrs() []string {
	addrs := []string{}

	for _, addr := range b.options.Nodes {
		if len(addr) == 0 {
			continue
		}

		if !strings.HasPrefix(addr, "nats://") {
			addr = "nats://" + addr
		}

		addrs = append(addrs, addr)
	}

	if len(addrs) == 0 {
		addrs = []string{client.DefaultURL}
	}

	return addrs
}

func (b *natsBroker) onClose(connection *client.Conn) {
	b.close <- nil
}

func (b *natsBroker) onAsyncError(connection *client.Conn, subscription *client.Subscription, err error) {
	if err == client.ErrDrainTimeout {
		b.close <- err
	}
}

func (b *natsBroker) onDisconnecedError(connection *client.Conn, err error) {
	b.close <- err
}

func NewBroker(opts ...broker.BrokerOption) broker.Broker {
	options := broker.NewBrokerOptions(opts...)

	b := &natsBroker{
		options: options,
		mtx:     sync.RWMutex{},
	}

	if err := b.configure(); err != nil {
		log.Fatal(err)
	}

	return b
}
