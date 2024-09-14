package grpcserver

import (
	"reflect"

	"github.com/w-h-a/pkg/server"
)

type grpcHandler struct {
	options  server.HandlerOptions
	name     string
	receiver reflect.Value
	methods  map[string]*grpcMethod
}

type grpcMethod struct {
	name    string
	value   reflect.Value
	ctxType reflect.Type
	reqType reflect.Type
	rspType reflect.Type
	stream  bool
}

func (c *grpcHandler) Options() server.HandlerOptions {
	return c.options
}

func (c *grpcHandler) Name() string {
	return c.name
}

func (c *grpcHandler) String() string {
	return "grpc"
}

func NewHandler(handler interface{}, opts ...server.HandlerOption) server.Handler {
	// TODO: handler options
	options := server.HandlerOptions{}

	methods := map[string]*grpcMethod{}

	// used to get method data
	typeOfHandler := reflect.TypeOf(handler)

	for i := 0; i < typeOfHandler.NumMethod(); i++ {
		m := typeOfHandler.Method(i)

		method := &grpcMethod{
			name:  m.Name,
			value: m.Func,
		}

		// TODO: make this better/safer
		switch m.Type.NumIn() {
		case 3:
			method.ctxType = m.Type.In(1)
			method.rspType = m.Type.In(2)
			method.stream = true
		case 4:
			method.ctxType = m.Type.In(1)
			method.reqType = m.Type.In(2)
			method.rspType = m.Type.In(3)
		}

		methods[m.Name] = method
	}

	// keep the value to use as a receiver in function invocation
	valueOfHandler := reflect.ValueOf(handler)

	// get the name of the handler struct
	nameOfHandler := reflect.Indirect(valueOfHandler).Type().Name()

	c := &grpcHandler{
		options:  options,
		name:     nameOfHandler,
		receiver: valueOfHandler,
		methods:  methods,
	}

	return c
}
