package api

import "context"

type ApiOption func(o *ApiOptions)

type ApiOptions struct {
	Namespace       string
	Name            string
	Version         string
	Address         string
	HandlerWrappers []HandlerWrapper
	Context         context.Context
}

func ApiWithNamespace(ns string) ApiOption {
	return func(o *ApiOptions) {
		o.Namespace = ns
	}
}

func ApiWithName(n string) ApiOption {
	return func(o *ApiOptions) {
		o.Name = n
	}
}

func ApiWithVersion(v string) ApiOption {
	return func(o *ApiOptions) {
		o.Version = v
	}
}

func ApiWithAddress(addr string) ApiOption {
	return func(o *ApiOptions) {
		o.Address = addr
	}
}

func WrapHandler(w HandlerWrapper) ApiOption {
	return func(o *ApiOptions) {
		o.HandlerWrappers = append(o.HandlerWrappers, w)
	}
}

func NewApiOptions(opts ...ApiOption) ApiOptions {
	options := ApiOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
