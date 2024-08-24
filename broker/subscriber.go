package broker

type Subscriber interface {
	Options() SubscribeOptions
	Id() string
	Handler(b []byte) error
	Unsubscribe() error
	String() string
}
