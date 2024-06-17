package api

import "net/http"

type Api interface {
	Options() ApiOptions
	Handle(path string, handler http.Handler)
	Run() error
	String() string
}
