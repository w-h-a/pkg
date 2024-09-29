package grpc

import (
	"fmt"
	"strconv"

	"github.com/w-h-a/pkg/runner"
	"github.com/w-h-a/pkg/serverv2"
	"github.com/w-h-a/pkg/serverv2/grpc"
	"github.com/w-h-a/pkg/telemetry/log"
)

type grpcProcess struct {
	options runner.ProcessOptions
	server  serverv2.Server
}

func (p *grpcProcess) Options() runner.ProcessOptions {
	return p.options
}

func (p *grpcProcess) Apply() error {
	return p.server.Run()
}

func (p *grpcProcess) Destroy() error {
	return p.server.Stop()
}

func (p *grpcProcess) String() string {
	return "grpc"
}

func NewProcess(opts ...runner.ProcessOption) runner.Process {
	options := runner.NewProcessOptions(opts...)

	var port int

	if prt, ok := options.EnvVars["PORT"]; ok {
		var err error
		port, err = strconv.Atoi(prt)
		if err != nil {
			log.Fatal(err)
		}
	}

	grpcServer := grpc.NewServer(
		serverv2.ServerWithNamespace("default"),
		serverv2.ServerWithName(options.Id),
		serverv2.ServerWithVersion("0.1.0"),
		serverv2.ServerWithAddress(fmt.Sprintf(":%d", port)),
	)

	if handlers, ok := GetHandlersFromContext(options.Context); ok {
		for _, handler := range handlers {
			if err := grpcServer.Handle(handler); err != nil {
				log.Fatal(err)
			}
		}
	}

	p := &grpcProcess{
		options: options,
		server:  grpcServer,
	}

	return p
}
