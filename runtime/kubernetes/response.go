package kubernetes

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type response struct {
	res *http.Response
	err error
}

func (r *response) decode(data interface{}) error {
	if r.err != nil {
		return r.err
	}

	defer r.res.Body.Close()

	decoder := json.NewDecoder(r.res.Body)

	if err := decoder.Decode(&data); err != nil {
		return ErrDecode
	}

	return r.err
}

func (r *response) getError() error {
	return r.err
}

func newResponse(res *http.Response, err error) *response {
	r := &response{
		res: res,
		err: err,
	}

	if err != nil {
		return r
	}

	if r.res.StatusCode == http.StatusOK {
		return r
	}

	if r.res.StatusCode == http.StatusNotFound {
		r.err = ErrResourceNotFound
		return r
	}

	b, err := io.ReadAll(r.res.Body)
	if err == nil {
		r.err = errors.New(string(b))
		return r
	}

	r.err = ErrUnknown

	return r
}
