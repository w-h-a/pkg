package broker

import "errors"

var (
	ErrNotConnected = errors.New("not connected")
)

type Broker interface {
	Options() BrokerOptions
	Publish(data interface{}, options PublishOptions) error
	Subscribe(callback func([]byte) error, options SubscribeOptions)
	String() string
}
