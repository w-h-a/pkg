package client

import (
	"context"
	"time"
)

type ClientOption func(o *ClientOptions)

type ClientOptions struct {
	ContentType    string
	CallOptions    CallOptions
	Selector       Selector
	ClientWrappers []ClientWrapper
	Context        context.Context
}

func ClientWithContentType(ct string) ClientOption {
	return func(o *ClientOptions) {
		o.ContentType = ct
	}
}

func ClientWithSelector(s Selector) ClientOption {
	return func(o *ClientOptions) {
		o.Selector = s
	}
}

func WrapClient(ws ...ClientWrapper) ClientOption {
	return func(o *ClientOptions) {
		o.ClientWrappers = append(o.ClientWrappers, ws...)
	}
}

func NewClientOptions(opts ...ClientOption) ClientOptions {
	options := ClientOptions{
		CallOptions: CallOptions{
			Backoff:        defaultBackoff,
			RetryCheck:     defaultRetryCheck,
			RetryCount:     defaultRetryCount,
			RequestTimeout: defaultRequestTimeout,
		},
		Context: context.Background(),
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

func NewRequestOptions(opts ...RequestOption) RequestOptions {
	options := RequestOptions{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type CallOption func(o *CallOptions)

type CallOptions struct {
	Address        string
	Backoff        func(ctx context.Context, req Request, attempts int) (time.Duration, error)
	RetryCheck     func(ctx context.Context, req Request, retryCount int, err error) (bool, error)
	RetryCount     int
	RequestTimeout time.Duration
	SelectOpts     []SelectOption
	CallWrappers   []CallWrapper
}

func CallWithAddress(addr string) CallOption {
	return func(o *CallOptions) {
		o.Address = addr
	}
}

func CallWithRetryCount(count int) CallOption {
	return func(o *CallOptions) {
		o.RetryCount = count
	}
}

func CallWithRequestTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.RequestTimeout = d
	}
}

func CallWithSelectOpts(opts ...SelectOption) CallOption {
	return func(o *CallOptions) {
		o.SelectOpts = append(o.SelectOpts, opts...)
	}
}

func WrapCall(ws ...CallWrapper) CallOption {
	return func(o *CallOptions) {
		o.CallWrappers = append(o.CallWrappers, ws...)
	}
}

func NewCallOptions(options *CallOptions, opts ...CallOption) CallOptions {
	for _, fn := range opts {
		fn(options)
	}

	return *options
}

type SelectorOption func(o *SelectorOptions)

type SelectorOptions struct {
	Context context.Context
}

func NewSelectorOptions(opts ...SelectorOption) SelectorOptions {
	options := SelectorOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type SelectOption func(o *SelectOptions)

type SelectOptions struct{}
