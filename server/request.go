package server

type Request interface {
	Options() RequestOptions
	Namespace() string
	Server() string
	Method() string
	ContentType() string
	Marshaled() ([]byte, error)
	Unmarshaled() interface{}
	String() string
}
