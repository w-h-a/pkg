package broker

type Publication interface {
	Topic() string
	Message() *Message
	Ack() error
	String() string
}
