package serverv2

import (
	"context"

	"github.com/w-h-a/pkg/telemetry/log"
)

type ServerOption func(o *ServerOptions)

type ServerOptions struct {
	Namespace string
	Name      string
	Id        string
	Version   string
	Address   string
	Logger    log.Log
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

func ServerWithLogger(l log.Log) ServerOption {
	return func(o *ServerOptions) {
		o.Logger = l
	}
}

func NewServerOptions(opts ...ServerOption) ServerOptions {
	options := ServerOptions{
		Namespace: defaultNamespace,
		Name:      defaultName,
		Id:        defaultID,
		Version:   defaultVersion,
		Address:   defaultAddress,
		Logger:    defaultLogger,
		Context:   context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
