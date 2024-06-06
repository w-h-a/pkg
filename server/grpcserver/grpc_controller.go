package grpcserver

import (
	"reflect"

	"github.com/w-h-a/pkg/server"
)

type grpcController struct {
	options  server.ControllerOptions
	name     string
	receiver reflect.Value
	handlers map[string]*grpcSync
}

type grpcSync struct {
	name    string
	method  reflect.Value
	ctxType reflect.Type
	reqType reflect.Type
	rspType reflect.Type
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

	handlers := map[string]*grpcSync{}

	// used to get method data
	typeOfController := reflect.TypeOf(controller)

	for i := 0; i < typeOfController.NumMethod(); i++ {
		method := typeOfController.Method(i)

		synchronousHandler := &grpcSync{
			name:   method.Name,
			method: method.Func,
		}

		// TODO: make this better
		switch method.Type.NumIn() {
		case 4:
			synchronousHandler.ctxType = method.Type.In(1)
			synchronousHandler.reqType = method.Type.In(2)
			synchronousHandler.rspType = method.Type.In(3)
		}

		handlers[method.Name] = synchronousHandler
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
