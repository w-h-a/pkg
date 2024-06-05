package main

import (
	"context"
	"log"

	proto "github.com/w-h-a/pkg/proto/runtime"
	"github.com/w-h-a/pkg/runtime"
	"github.com/w-h-a/pkg/runtime/kubernetes"
	"github.com/w-h-a/pkg/server"
	"github.com/w-h-a/pkg/server/grpcserver"
)

func main() {
	grpcServer := grpcserver.NewServer(
		server.ServerWithNamespace("default"),
		server.ServerWithName("runtime"),
		server.ServerWithVersion("v1.0.0"),
		server.ServerWithAddress(":8080"),
	)

	r := kubernetes.NewRuntime()

	if err := RegisterRuntimeController(grpcServer, NewRuntimeController(r)); err != nil {
		log.Fatalf("failed to register controller: %v", err)
	}

	if err := grpcServer.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

type RuntimeController interface {
	Get(ctx context.Context, req *proto.GetRequest, rsp *proto.GetResponse) error
}

type runtimeController struct {
	runtime runtime.Runtime
}

func (c *runtimeController) Get(ctx context.Context, req *proto.GetRequest, rsp *proto.GetResponse) error {
	services, err := c.runtime.GetServices()
	if err != nil {
		return err
	}

	result := []*proto.Service{}

	for _, service := range services {
		proto := &proto.Service{
			Namespace: service.Namespace,
			Name: service.Name,
			Version: service.Version,
			Address: service.Address,
			Metadata: service.Metadata,
		}
		result = append(result, proto)
	}

	rsp.Services = result

	return nil
}

func NewRuntimeController(r runtime.Runtime) RuntimeController {
	return &runtimeController{r}
}

type Runtime struct {
	RuntimeController
}

func RegisterRuntimeController(s server.Server, controller RuntimeController, opts ...server.ControllerOption) error {
	return s.RegisterController(
		s.NewController(
			&Runtime{controller},
			opts...,
		),
	)
}