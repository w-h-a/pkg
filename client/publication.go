package client

type Publication interface {
	Options() PublicationOptions
	Topic() string
	ContentType() string
	Unmarshaled() interface{}
	String() string
}
