package server

type Request interface {
	Options() RequestOptions
	Namespace() string
	Service() string
	Method() string
	ContentType() string
	Unmarshaled() interface{}
	Marshaled() []byte
	Stream() bool
	String() string
}
