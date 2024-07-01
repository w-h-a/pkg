package streams

var ()

type Stream interface {
	Produce(topic string, payload interface{}, opts ...ProduceOption) error
	Consume(topic string, opts ...ConsumeOption) (<-chan Event, error)
}
