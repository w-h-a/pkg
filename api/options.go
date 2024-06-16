package api

import "context"

type ApiOption func(o *ApiOptions)

type ApiOptions struct {
	Address         string
	HandlerWrappers []HandlerWrapper
	Context         context.Context
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
