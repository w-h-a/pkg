package server

type Handler interface {
	Options() HandlerOptions
	Name() string
	String() string
}
