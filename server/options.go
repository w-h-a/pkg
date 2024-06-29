package server

import (
	"context"

	"github.com/w-h-a/pkg/broker"
)

type ServerOption func(o *ServerOptions)

type ServerOptions struct {
	Namespace          string
	Name               string
	Id                 string
	Version            string
	Address            string
	Metadata           map[string]string
	Broker             broker.Broker
	ControllerWrappers []ControllerWrapper
	SubscriberWrappers []SubscriberWrapper
	Context            context.Context
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

func ServerWithMetadata(md map[string]string) ServerOption {
	return func(o *ServerOptions) {
		o.Metadata = md
	}
}

func ServerWithBroker(b broker.Broker) ServerOption {
	return func(o *ServerOptions) {
		o.Broker = b
	}
}

func WrapController(ws ...ControllerWrapper) ServerOption {
	return func(o *ServerOptions) {
		o.ControllerWrappers = append(o.ControllerWrappers, ws...)
	}
}

func WrapSubscriber(ws ...SubscriberWrapper) ServerOption {
	return func(o *ServerOptions) {
		o.SubscriberWrappers = append(o.SubscriberWrappers, ws...)
	}
}

func NewServerOptions(opts ...ServerOption) ServerOptions {
	options := ServerOptions{
		Namespace: defaultNamespace,
		Name:      defaultName,
		Id:        defaultID,
		Version:   defaultVersion,
		Address:   defaultAddress,
		Metadata:  map[string]string{},
		Context:   context.Background(),
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
	MarshaledRequest   []byte
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

func RequestWithMarshaledRequest(bs []byte) RequestOption {
	return func(o *RequestOptions) {
		o.MarshaledRequest = bs
	}
}

func NewRequestOptions(opts ...RequestOption) RequestOptions {
	options := RequestOptions{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type PublicationOption func(o *PublicationOptions)

type PublicationOptions struct {
	Topic              string
	ContentType        string
	UnmarshaledPayload interface{}
}

func PublicationWithTopic(t string) PublicationOption {
	return func(o *PublicationOptions) {
		o.Topic = t
	}
}

func PublicationWithContentType(ct string) PublicationOption {
	return func(o *PublicationOptions) {
		o.ContentType = ct
	}
}

func PublicationWithUnmarshaledPayload(v interface{}) PublicationOption {
	return func(o *PublicationOptions) {
		o.UnmarshaledPayload = v
	}
}

func NewPublicationOptions(opts ...PublicationOption) PublicationOptions {
	options := PublicationOptions{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type ControllerOption func(o *ControllerOptions)

type ControllerOptions struct{}

type SubscriberOption func(o *SubscriberOptions)

type SubscriberOptions struct {
	AutoAck   bool
	QueueName string
}

func SubscriberWithoutAutoAck() SubscriberOption {
	return func(o *SubscriberOptions) {
		o.AutoAck = false
	}
}

func SubscriberWithQueueName(n string) SubscriberOption {
	return func(o *SubscriberOptions) {
		o.QueueName = n
	}
}

func NewSubscriberOptions(opts ...SubscriberOption) SubscriberOptions {
	options := SubscriberOptions{
		AutoAck: true,
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
