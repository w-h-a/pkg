package client

type Request interface {
	Options() RequestOptions
	Namespace() string
	Server() string
	Port() int
	Method() string
	ContentType() string
	Unmarshaled() interface{}
	String() string
}
