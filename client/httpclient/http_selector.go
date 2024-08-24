package httpclient

import (
	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/runtime"
)

type httpSelector struct {
	options client.SelectorOptions
}

func (s *httpSelector) Options() client.SelectorOptions {
	return s.options
}

func (s *httpSelector) Select(namespace, service string, port int, opts ...client.SelectOption) (func() (*runtime.Service, error), error) {
	return func() (*runtime.Service, error) {
		return &runtime.Service{
			Namespace: namespace,
			Name:      service,
			Port:      port,
		}, nil
	}, nil
}

func (s *httpSelector) String() string {
	return "http"
}

func NewSelector(opts ...client.SelectorOption) client.Selector {
	options := client.NewSelectorOptions(opts...)

	s := &httpSelector{
		options: options,
	}

	return s
}
