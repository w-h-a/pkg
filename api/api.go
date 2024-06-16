package api

import "net/http"

type Api interface {
	Options() ApiOptions
	Handle(path string, handler http.Handler)
	Start() error
	Stop() error
	String() string
}
