package handlers

import (
	"context"

	"github.com/w-h-a/pkg/examples/greeter/proto"
	"github.com/w-h-a/pkg/server"
	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/utils/metadatautils"
)

type GreetHandler interface {
	Greet(ctx context.Context, req *proto.GreetRequest, rsp *proto.GreetResponse) error
}

type greetHandler struct {}

func (c *greetHandler) Greet(ctx context.Context, req *proto.GreetRequest, rsp *proto.GreetResponse) error {
	md, _ := metadatautils.FromContext(ctx)
	log.Infof("received Greeter.Greet request with metadata: %v", md)

	rsp.Msg = "Hello, " + req.Name

	return nil
}

func NewGreetHandler() GreetHandler {
	return &greetHandler{}
}

type Greeter struct {
	GreetHandler
}

func RegisterGreetHandler(s server.Server, controller GreetHandler, opts ...server.HandlerOption) error {
	return s.Handle(
		s.NewHandler(
			&Greeter{controller},
			opts...,
		),
	)
}