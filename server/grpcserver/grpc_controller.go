package grpcserver

import (
	"reflect"

	"github.com/w-h-a/pkg/server"
)

type grpcController struct {
	options  server.ControllerOptions
	name     string
	receiver reflect.Value
	handlers map[string]*grpcHandler
}

type grpcHandler struct {
	name    string
	method  reflect.Value
	ctxType reflect.Type
	reqType reflect.Type
	rspType reflect.Type
	stream  bool
}

func (c *grpcController) Options() server.ControllerOptions {
	return c.options
}

func (c *grpcController) Name() string {
	return c.name
}

func (c *grpcController) String() string {
	return "grpc"
}

func NewController(controller interface{}, opts ...server.ControllerOption) server.Controller {
	// TODO: controller options
	options := server.ControllerOptions{}

	handlers := map[string]*grpcHandler{}

	// used to get method data
	typeOfController := reflect.TypeOf(controller)

	for i := 0; i < typeOfController.NumMethod(); i++ {
		method := typeOfController.Method(i)

		handler := &grpcHandler{
			name:   method.Name,
			method: method.Func,
		}

		// TODO: make this better/safer
		switch method.Type.NumIn() {
		case 3:
			handler.ctxType = method.Type.In(1)
			handler.rspType = method.Type.In(2)
			handler.stream = true
		case 4:
			handler.ctxType = method.Type.In(1)
			handler.reqType = method.Type.In(2)
			handler.rspType = method.Type.In(3)
		}

		handlers[method.Name] = handler
	}

	// keep the value to use as a receiver in function invocation
	valueOfController := reflect.ValueOf(controller)

	// get the name of the controller struct
	nameOfController := reflect.Indirect(valueOfController).Type().Name()

	c := &grpcController{
		options:  options,
		name:     nameOfController,
		receiver: valueOfController,
		handlers: handlers,
	}

	return c
}
