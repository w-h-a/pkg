package serverv2

import (
	"context"
)

type ServerOption func(o *ServerOptions)

type ServerOptions struct {
	Namespace string
	Name      string
	Id        string
	Version   string
	Address   string
	Tracer    string
	Context   context.Context
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

func ServerWithTracer(tracer string) ServerOption {
	return func(o *ServerOptions) {
		o.Tracer = tracer
	}
}

func NewServerOptions(opts ...ServerOption) ServerOptions {
	options := ServerOptions{
		Namespace: defaultNamespace,
		Name:      defaultName,
		Id:        defaultID,
		Version:   defaultVersion,
		Address:   defaultAddress,
		Context:   context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
