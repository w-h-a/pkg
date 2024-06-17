package grpcclient

import (
	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/runtime"
)

type grpcSelector struct {
	options client.SelectorOptions
}

func (s *grpcSelector) Options() client.SelectorOptions {
	return s.options
}

func (s *grpcSelector) Select(namespace, service string, port int, opts ...client.SelectOption) (func() (*runtime.Service, error), error) {
	return func() (*runtime.Service, error) {
		return &runtime.Service{
			Namespace: namespace,
			Name:      service,
			Port:      port,
		}, nil
	}, nil
}

func (s *grpcSelector) String() string {
	return "grpc"
}

func NewSelector(opts ...client.SelectorOption) client.Selector {
	options := client.NewSelectorOptions(opts...)

	s := &grpcSelector{
		options: options,
	}

	return s
}
