package server

import "github.com/google/uuid"

var (
	defaultNamespace = "default"
	defaultName      = "server"
	defaultID        = uuid.New().String()
	defaultVersion   = "v0.1.0"
	defaultAddress   = ":0"
)

type Server interface {
	Options() ServerOptions
	NewHandler(c interface{}, opts ...HandlerOption) Handler
	Handle(c Handler) error
	Run() error
	String() string
}
