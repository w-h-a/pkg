package server

type Subscriber interface {
	Options() SubscriberOptions
	Topic() string
	String() string
}
