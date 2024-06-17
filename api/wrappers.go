package api

import "net/http"

type HandlerWrapper func(h http.Handler) http.Handler
