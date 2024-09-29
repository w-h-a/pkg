package serverv2

import (
	"github.com/google/uuid"
	"github.com/w-h-a/pkg/telemetry/log/memory"
)

var (
	defaultNamespace = "default"
	defaultName      = "server"
	defaultID        = uuid.New().String()
	defaultVersion   = "v0.1.0"
	defaultAddress   = ":0"
	defaultLogger    = memory.NewLog()
)

type Server interface {
	Options() ServerOptions
	Handle(handler interface{}) error
	Start() error
	Run() error
	Stop() error
	String() string
}
