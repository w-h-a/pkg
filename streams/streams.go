package streams

import "errors"

var (
	ErrSubscriberNotFound = errors.New("failed to find subscriber")
	ErrEncodingData       = errors.New("failed to encode incoming data")
	ErrEncodingEvent      = errors.New("failed to encode outgoing event")
)

type Streams interface {
	Options() StreamsOptions
	Subscribe(id string, opts ...SubscribeOption) error
	Unsubscribe(id string) error
	Consume(id string) (Subscriber, error)
	Produce(topic string, data []byte, opts ...ProduceOption) error
	String() string
}
