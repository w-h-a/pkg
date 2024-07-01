package server

import "context"

type ServerOption func(o *ServerOptions)

type ServerOptions struct {
	Namespace          string
	Name               string
	Id                 string
	Version            string
	Address            string
	Metadata           map[string]string
	ControllerWrappers []ControllerWrapper
	Context            context.Context
}

func ServerWithNamespace(n string) ServerOption {
	return func(o *ServerOptions) {
		o.Namespace = n
	}
}

func ServerWithName(n string) ServerOption {
	return func(o *ServerOptions) {
		o.Name = n
	}
}

func ServerWithId(id string) ServerOption {
	return func(o *ServerOptions) {
		o.Id = id
	}
}

func ServerWithVersion(v string) ServerOption {
	return func(o *ServerOptions) {
		o.Version = v
	}
}

func ServerWithAddress(addr string) ServerOption {
	return func(o *ServerOptions) {
		o.Address = addr
	}
}

func ServerWithMetadata(md map[string]string) ServerOption {
	return func(o *ServerOptions) {
		o.Metadata = md
	}
}

func WrapController(ws ...ControllerWrapper) ServerOption {
	return func(o *ServerOptions) {
		o.ControllerWrappers = append(o.ControllerWrappers, ws...)
	}
}

func NewServerOptions(opts ...ServerOption) ServerOptions {
	options := ServerOptions{
		Namespace: defaultNamespace,
		Name:      defaultName,
		Id:        defaultID,
		Version:   defaultVersion,
		Address:   defaultAddress,
		Metadata:  map[string]string{},
		Context:   context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type RequestOption func(o *RequestOptions)

type RequestOptions struct {
	Namespace          string
	Name               string
	Method             string
	ContentType        string
	UnmarshaledRequest interface{}
	MarshaledRequest   []byte
	Stream             bool
}

func RequestWithNamespace(n string) RequestOption {
	return func(o *RequestOptions) {
		o.Namespace = n
	}
}

func RequestWithName(n string) RequestOption {
	return func(o *RequestOptions) {
		o.Name = n
	}
}

func RequestWithMethod(m string) RequestOption {
	return func(o *RequestOptions) {
		o.Method = m
	}
}

func RequestWithContentType(ct string) RequestOption {
	return func(o *RequestOptions) {
		o.ContentType = ct
	}
}

func RequestWithUnmarshaledRequest(v interface{}) RequestOption {
	return func(o *RequestOptions) {
		o.UnmarshaledRequest = v
	}
}

func RequestWithMarshaledRequest(bs []byte) RequestOption {
	return func(o *RequestOptions) {
		o.MarshaledRequest = bs
	}
}

func RequestWithStream() RequestOption {
	return func(o *RequestOptions) {
		o.Stream = true
	}
}

func NewRequestOptions(opts ...RequestOption) RequestOptions {
	options := RequestOptions{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type ControllerOption func(o *ControllerOptions)

type ControllerOptions struct{}
