package sidecar

import (
	"context"

	"github.com/w-h-a/pkg/broker"
	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/security/secret"
	"github.com/w-h-a/pkg/store"
	"github.com/w-h-a/pkg/telemetry/tracev2"
)

type SidecarOption func(o *SidecarOptions)

type SidecarOptions struct {
	ServiceName string
	HttpPort    Port
	GrpcPort    Port
	ServicePort Port
	Client      client.Client
	Stores      map[string]store.Store
	Brokers     map[string]broker.Broker
	Secrets     map[string]secret.Secret
	Tracer      tracev2.Trace
	Context     context.Context
}

type Port struct {
	Port     string
	Protocol string
}

func SidecarWithServiceName(n string) SidecarOption {
	return func(o *SidecarOptions) {
		o.ServiceName = n
	}
}

func SidecarWithHttpPort(p Port) SidecarOption {
	return func(o *SidecarOptions) {
		o.HttpPort = p
	}
}

func SidecarWithGrpcPort(p Port) SidecarOption {
	return func(o *SidecarOptions) {
		o.GrpcPort = p
	}
}

func SidecarWithServicePort(p Port) SidecarOption {
	return func(o *SidecarOptions) {
		o.ServicePort = p
	}
}

func SidecarWithClient(c client.Client) SidecarOption {
	return func(o *SidecarOptions) {
		o.Client = c
	}
}

func SidecarWithStores(s map[string]store.Store) SidecarOption {
	return func(o *SidecarOptions) {
		o.Stores = s
	}
}

func SidecarWithBrokers(b map[string]broker.Broker) SidecarOption {
	return func(o *SidecarOptions) {
		o.Brokers = b
	}
}

func SidecarWithSecrets(s map[string]secret.Secret) SidecarOption {
	return func(o *SidecarOptions) {
		o.Secrets = s
	}
}

func SidecarWithTracer(tr tracev2.Trace) SidecarOption {
	return func(o *SidecarOptions) {
		o.Tracer = tr
	}
}

func NewSidecarOptions(opts ...SidecarOption) SidecarOptions {
	options := SidecarOptions{
		Stores:  map[string]store.Store{},
		Brokers: map[string]broker.Broker{},
		Secrets: map[string]secret.Secret{},
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
