package server

type Request interface {
	Options() RequestOptions
	Namespace() string
	Server() string
	Method() string
	ContentType() string
	Unmarshaled() interface{}
	Marshaled() []byte
	String() string
}
