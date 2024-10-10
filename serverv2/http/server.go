package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/w-h-a/pkg/serverv2"
	"github.com/w-h-a/pkg/serverv2/http/handlers"
	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/telemetry/trace"
)

type server struct {
	options serverv2.ServerOptions
	mux     *http.ServeMux
	started bool
	mtx     sync.RWMutex
	errCh   chan error
	exit    chan struct{}
}

func (s *server) Options() serverv2.ServerOptions {
	return s.options
}

func (s *server) Handle(handler interface{}) error {
	h, ok := handler.(http.Handler)
	if !ok {
		return fmt.Errorf("invalid handler: expected http.Handler")
	}

	if ms, ok := GetMiddlewaresFromContext(s.options.Context); ok && ms != nil {
		for i := len(ms); i > 0; i-- {
			h = ms[i-1](h)
		}
	}

	s.mux.Handle("/", h)

	return nil
}

func (s *server) Start() error {
	if err := s.Run(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	log.Infof("http server received signal %s", <-ch)

	return s.Stop()
}

func (s *server) Run() error {
	s.mtx.RLock()
	if s.started {
		s.mtx.RUnlock()
		return nil
	}
	s.mtx.RUnlock()

	if len(s.options.Tracer) > 0 {
		// init trace exporters
		switch s.options.Tracer {
		case "memory":
			tracer := trace.GetTracer()
			if tracer == nil {
				log.Fatalf("failed to init memory trace exporter: memory tracer was not set")
			}
			router := mux.NewRouter()
			httpTrace := handlers.NewTraceHandler(tracer)
			router.Methods("GET").Path("/").HandlerFunc(httpTrace.Read)
			s.mux.Handle("/trace", router)
		default:
			log.Warnf("tracer %s is not supported", s.options.Tracer)
		}
	}

	// TODO: tls
	listener, err := net.Listen("tcp", s.options.Address)
	if err != nil {
		return err
	}

	s.mtx.Lock()
	s.options.Address = listener.Addr().String()
	s.mtx.Unlock()

	log.Infof("http server is listening on %s", s.options.Address)

	httpServer := &http.Server{
		ReadHeaderTimeout: time.Second,
		Handler:           s.mux,
		TLSConfig:         nil,
	}

	go func() {
		s.errCh <- httpServer.Serve(listener)
	}()

	go func() {
		<-s.exit

		var err error

		shutdown := make(chan struct{})

		go func() {
			defer close(shutdown)
			err = httpServer.Shutdown(context.Background())
		}()

		select {
		case <-shutdown:
		case <-time.After(10 * time.Second):
			err = httpServer.Close()
		}

		s.errCh <- err
	}()

	s.mtx.Lock()
	s.started = true
	s.mtx.Unlock()

	return nil
}

func (s *server) Stop() error {
	s.mtx.RLock()
	if !s.started {
		s.mtx.RUnlock()
		return nil
	}
	s.mtx.RUnlock()

	close(s.exit)

	var err error

	for i := 0; i < 2; i++ {
		err = <-s.errCh
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
	}

	s.mtx.Lock()
	s.started = false
	s.mtx.Unlock()

	return err
}

func (s *server) String() string {
	return "http"
}

func NewServer(opts ...serverv2.ServerOption) serverv2.Server {
	options := serverv2.NewServerOptions(opts...)

	s := &server{
		options: options,
		mux:     http.NewServeMux(),
		mtx:     sync.RWMutex{},
		errCh:   make(chan error, 2),
		exit:    make(chan struct{}),
	}

	return s
}
