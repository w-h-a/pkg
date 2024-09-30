package serverv2

import (
	"github.com/google/uuid"
)

var (
	defaultNamespace = "default"
	defaultName      = "server"
	defaultID        = uuid.New().String()
	defaultVersion   = "v0.1.0"
	defaultAddress   = ":0"
)

type Server interface {
	Options() ServerOptions
	Handle(handler interface{}) error
	Start() error
	Run() error
	Stop() error
	String() string
}
