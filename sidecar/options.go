package sidecar

import (
	"context"

	"github.com/w-h-a/pkg/broker"
	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/store"
)

type SidecarOption func(o *SidecarOptions)

type SidecarOptions struct {
	Id          string
	ServiceName string
	HttpPort    Port
	RpcPort     Port
	ServicePort Port
	HttpClient  client.Client
	RpcClient   client.Client
	Stores      map[string]store.Store
	Brokers     map[string]broker.Broker
	Context     context.Context
}

type Port struct {
	Port     string
	Protocol string
}

func SidecarWithId(id string) SidecarOption {
	return func(o *SidecarOptions) {
		o.Id = id
	}
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

func SidecarWithRpcPort(p Port) SidecarOption {
	return func(o *SidecarOptions) {
		o.RpcPort = p
	}
}

func SidecarWithServicePort(p Port) SidecarOption {
	return func(o *SidecarOptions) {
		o.ServicePort = p
	}
}

func SidecarWithHttpClient(c client.Client) SidecarOption {
	return func(o *SidecarOptions) {
		o.HttpClient = c
	}
}

func SidecarWithRpcClient(c client.Client) SidecarOption {
	return func(o *SidecarOptions) {
		o.RpcClient = c
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

func NewSidecarOptions(opts ...SidecarOption) SidecarOptions {
	options := SidecarOptions{
		Id:      defaultID,
		Stores:  map[string]store.Store{},
		Brokers: map[string]broker.Broker{},
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
