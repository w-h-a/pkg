package broker

import "errors"

var (
	ErrNotConnected = errors.New("not connected")
)

type Broker interface {
	Options() BrokerOptions
	Address() string
	Connect() error
	Disconnect() error
	Subscribe(topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error)
	Publish(topic string, msg *Message) error
	String() string
}
