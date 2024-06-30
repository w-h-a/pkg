package broker

type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe() error
	String() string
}
