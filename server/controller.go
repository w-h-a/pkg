package server

type Controller interface {
	Options() ControllerOptions
	Name() string
	String() string
}
