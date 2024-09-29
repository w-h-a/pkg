package http

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/w-h-a/pkg/runner"
	"github.com/w-h-a/pkg/serverv2"
	"github.com/w-h-a/pkg/serverv2/http"
)

type httpProcess struct {
	options runner.ProcessOptions
	server  serverv2.Server
}

func (p *httpProcess) Options() runner.ProcessOptions {
	return p.options
}

func (p *httpProcess) Apply() error {
	return p.server.Run()
}

func (p *httpProcess) Destroy() error {
	return p.server.Stop()
}

func (p *httpProcess) String() string {
	return "http"
}

func NewProcess(opts ...runner.ProcessOption) runner.Process {
	options := runner.NewProcessOptions(opts...)

	router := mux.NewRouter()

	if handlers, ok := GetHandlersFromContext(options.Context); ok {
		for path, handler := range handlers {
			router.Path(path).HandlerFunc(handler)
		}
	}

	var port int

	if prt, ok := options.EnvVars["PORT"]; ok {
		var err error
		port, err = strconv.Atoi(prt)
		if err != nil {
			log.Fatal(err)
		}
	}

	httpServer := http.NewServer(
		serverv2.ServerWithNamespace("default"),
		serverv2.ServerWithName(options.Id),
		serverv2.ServerWithVersion("0.1.0"),
		serverv2.ServerWithAddress(fmt.Sprintf(":%d", port)),
	)

	if err := httpServer.Handle(router); err != nil {
		log.Fatal(err)
	}

	p := &httpProcess{
		options: options,
		server:  httpServer,
	}

	return p
}
