package api

import (
	"net/http"

	"github.com/google/uuid"
)

var (
	defaultNamespace = "default"
	defaultName      = "api"
	defaultID        = uuid.New().String()
	defaultVersion   = "v0.1.0"
	defaultAddress   = ":0"
)

type Api interface {
	Options() ApiOptions
	Handle(path string, handler http.Handler)
	Run() error
	String() string
}
