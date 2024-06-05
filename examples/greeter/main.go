package main

import (
	"context"
	"strings"

	"github.com/w-h-a/pkg/examples/greeter/controllers"
	"github.com/w-h-a/pkg/server"
	"github.com/w-h-a/pkg/server/grpcserver"
	"github.com/w-h-a/pkg/telemetry/log"
)

func main() {
	grpcServer := grpcserver.NewServer(
		server.ServerWithNamespace("app"),
		server.ServerWithName("greeter"),
		server.ServerWithVersion("v1.0.0"),
		server.WrapController(controllerLogWrapper),
	)

	if err := controllers.RegisterGreetController(grpcServer, controllers.NewGreetController()); err != nil {
		log.Fatalf("failed to register controller: %v", err)
	}

	if err := grpcServer.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func controllerLogWrapper(fn server.ControllerFunc) server.ControllerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		if strings.HasPrefix(req.Method(), "Health.") {
			return fn(ctx, req, rsp)
		}

		log.Infof("before serving request for: %v", req.Method())

		if err := fn(ctx, req, rsp); err != nil {
			log.Errorf("method %v failed: %v", req.Method(), err)
			return err
		}

		log.Infof("after serving request for: %v", req.Method())

		return nil
	}
}