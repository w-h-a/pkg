package streams

type Subscriber interface {
	Options() SubscribeOptions
	Channel() chan Event
	Close()
	Ack(ev Event) interface{}
	Nack(ev Event) interface{}
	SetAttempts(c int, ev Event)
	GetAttempts(ev Event) (int, bool)
	String() string
}
