package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/w-h-a/pkg/runner"
	"github.com/w-h-a/pkg/telemetry/log"
)

type httpProcess struct {
	options  runner.ProcessOptions
	listener net.Listener
	server   *http.Server
	errCh    chan error
	exit     chan struct{}
}

func (p *httpProcess) Options() runner.ProcessOptions {
	return p.options
}

func (p *httpProcess) Apply() error {
	go func() {
		err := p.server.Serve(p.listener)
		if !errors.Is(err, http.ErrServerClosed) {
			p.errCh <- err
		} else {
			p.errCh <- nil
		}
	}()

	go func() {
		<-p.exit
		log.Infof("http process received interrupt signal")
		p.errCh <- p.server.Shutdown(context.Background())
	}()

	return nil
}

func (p *httpProcess) Destroy() error {
	close(p.exit)

	var err error

	for i := 0; i < 2; i++ {
		err = <-p.errCh
	}

	return err
}

func (p *httpProcess) String() string {
	return "http"
}

func NewProcess(opts ...runner.ProcessOption) runner.Process {
	options := runner.NewProcessOptions(opts...)

	mux := http.NewServeMux()

	if handlers, ok := GetHandlerFuncsFromContext(options.Context); ok {
		for path, handler := range handlers {
			mux.HandleFunc(path, handler)
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

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("http service is listening on port %d", port)

	p := &httpProcess{
		options:  options,
		listener: listener,
		server: &http.Server{
			ReadHeaderTimeout: time.Second,
			Handler:           mux,
			TLSConfig:         nil,
		},
		errCh: make(chan error, 2),
		exit:  make(chan struct{}),
	}

	return p
}