package httpclient

import "github.com/w-h-a/pkg/client"

type httpRequest struct {
	options client.RequestOptions
}

func (r *httpRequest) Options() client.RequestOptions {
	return r.options
}

func (r *httpRequest) Namespace() string {
	return r.options.Namespace
}

func (r *httpRequest) Service() string {
	return r.options.Name
}

func (r *httpRequest) Method() string {
	return r.options.Method
}

func (r *httpRequest) Port() int {
	return r.options.Port
}

func (r *httpRequest) ContentType() string {
	return r.options.ContentType
}

func (r *httpRequest) Unmarshaled() interface{} {
	return r.options.UnmarshaledRequest
}

func (r *httpRequest) Stream() bool {
	return r.options.Stream
}

func (r *httpRequest) String() string {
	return "http"
}

func NewRequest(opts ...client.RequestOption) client.Request {
	options := client.NewRequestOptions(opts...)

	if len(options.ContentType) == 0 {
		options.ContentType = defaultContentType
	}

	return &httpRequest{
		options: options,
	}
}
