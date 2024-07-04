package streams

type Subscriber interface {
	Options() SubscribeOptions
	Channel() chan Event
	Ack(ev Event) error
	Nack(ev Event) error
	SetAttempts(c int, ev Event)
	GetAttempts(ev Event) (int, bool)
	String() string
}
