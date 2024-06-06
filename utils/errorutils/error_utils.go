package errorutils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Id     string `json:"id"`
	Code   int32  `json:"code"`
	Detail string `json:"detail"`
	Status string `json:"status"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func NewError(id, detail string, code int32) error {
	return &Error{
		Id:     id,
		Code:   code,
		Detail: detail,
		Status: http.StatusText(int(code)),
	}
}

func ParseError(errStr string) *Error {
	e := &Error{}

	if err := json.Unmarshal([]byte(errStr), e); err != nil {
		e.Detail = errStr
	}

	return e
}

func Timeout(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   408,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(http.StatusRequestTimeout),
	}
}

func InternalServerError(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   500,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(http.StatusInternalServerError),
	}
}
