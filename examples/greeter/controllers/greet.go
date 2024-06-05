package controllers

import (
	"context"

	"github.com/w-h-a/pkg/examples/greeter/proto"
	"github.com/w-h-a/pkg/server"
	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/utils/metadatautils"
)

type GreetController interface {
	Greet(ctx context.Context, req *proto.GreetRequest, rsp *proto.GreetResponse) error
}

type greetController struct {}

func (c *greetController) Greet(ctx context.Context, req *proto.GreetRequest, rsp *proto.GreetResponse) error {
	md, _ := metadatautils.FromContext(ctx)
	log.Infof("received Greeter.Greet request with metadata: %v", md)

	rsp.Msg = "Hello, " + req.Name

	return nil
}

func NewGreetController() GreetController {
	return &greetController{}
}

type Greeter struct {
	GreetController
}

func RegisterGreetController(s server.Server, controller GreetController, opts ...server.ControllerOption) error {
	return s.RegisterController(
		s.NewController(
			&Greeter{controller},
			opts...,
		),
	)
}