package broker

import (
	"context"
	"crypto/tls"

	"github.com/w-h-a/pkg/utils/marshalutils"
)

type BrokerOption func(o *BrokerOptions)

type BrokerOptions struct {
	Nodes           []string
	Secure          bool
	TLSConfig       *tls.Config
	Marshaler       marshalutils.Marshaler
	HandlerWrappers []HandlerWrapper
	Context         context.Context
}

func BrokerWithNodes(addrs ...string) BrokerOption {
	return func(o *BrokerOptions) {
		o.Nodes = addrs
	}
}

func BrokerWithSecure() BrokerOption {
	return func(o *BrokerOptions) {
		o.Secure = true
	}
}

func BrokerWithTLSConfig(cfg *tls.Config) BrokerOption {
	return func(o *BrokerOptions) {
		o.TLSConfig = cfg
	}
}

func BrokerWithMarshaler(m marshalutils.Marshaler) BrokerOption {
	return func(o *BrokerOptions) {
		o.Marshaler = m
	}
}

func WrapHandler(ws ...HandlerWrapper) BrokerOption {
	return func(o *BrokerOptions) {
		o.HandlerWrappers = append(o.HandlerWrappers, ws...)
	}
}

func NewBrokerOptions(opts ...BrokerOption) BrokerOptions {
	options := BrokerOptions{
		Marshaler: marshalutils.DefaultMarshalers["application/json"],
		Context:   context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type SubscribeOption func(o *SubscribeOptions)

type SubscribeOptions struct {
	AutoAck bool
	Queue   string
}

func SubscribeWithoutAutoAck() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = false
	}
}

func SubscribeWithQueue(n string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = n
	}
}

func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	options := SubscribeOptions{
		AutoAck: true,
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
