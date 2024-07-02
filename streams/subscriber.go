package streams

type Subscriber interface {
	Options() SubscribeOptions
	Channel() chan Event
	Ack(ev Event) error
	Nack(ev Event) error
	SetAttemptCount(c int, ev Event)
	GetAttemptCount(ev Event) (int, bool)
	String() string
}
