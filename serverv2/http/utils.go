package http

import (
	"encoding/json"
	"net/http"

	"github.com/w-h-a/pkg/utils/errorutils"
)

func ErrResponse(w http.ResponseWriter, err error) {
	internal := err.(*errorutils.Error)
	Response(w, int(internal.Code), []byte(internal.Error()))
}

func OkResponse(w http.ResponseWriter, payload interface{}) {
	bs, _ := json.Marshal(payload)
	Response(w, 200, bs)
}

func Response(w http.ResponseWriter, code int, bs []byte) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(bs)
}