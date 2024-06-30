package grpcserver

import (
	"reflect"

	"github.com/w-h-a/pkg/server"
)

type grpcSubscriber struct {
	options  server.SubscriberOptions
	topic    string
	receiver reflect.Value
	handlers []*grpcAsync
}

type grpcAsync struct {
	name        string
	method      reflect.Value
	ctxType     reflect.Type
	payloadType reflect.Type
}

func (s *grpcSubscriber) Options() server.SubscriberOptions {
	return s.options
}

func (s *grpcSubscriber) Topic() string {
	return s.topic
}

func (s *grpcSubscriber) String() string {
	return "grpc"
}

func NewSubscriber(topic string, subscriber interface{}, opts ...server.SubscriberOption) server.Subscriber {
	options := server.NewSubscriberOptions(opts...)

	handlers := []*grpcAsync{}

	// used to get method data
	typeOfSubscriber := reflect.TypeOf(subscriber)

	for i := 0; i < typeOfSubscriber.NumMethod(); i++ {
		method := typeOfSubscriber.Method(i)

		asynchronousHandler := &grpcAsync{
			name:   method.Name,
			method: method.Func,
		}

		// TODO: make this better
		switch method.Type.NumIn() {
		case 3:
			asynchronousHandler.ctxType = method.Type.In(1)
			asynchronousHandler.payloadType = method.Type.In(2)
		}

		handlers = append(handlers, asynchronousHandler)
	}

	// keep the value to use as receiver in function invocation
	valueOfSubscriber := reflect.ValueOf(subscriber)

	s := &grpcSubscriber{
		options:  options,
		topic:    topic,
		receiver: valueOfSubscriber,
		handlers: handlers,
	}

	return s
}
