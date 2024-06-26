package httpapi

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/w-h-a/pkg/api"
	"github.com/w-h-a/pkg/telemetry/log"
)

type httpApi struct {
	options api.ApiOptions
	mux     *http.ServeMux
	mtx     sync.RWMutex
	exit    chan chan error
}

func (a *httpApi) Options() api.ApiOptions {
	return a.options
}

func (a *httpApi) Handle(path string, handler http.Handler) {
	h := handlers.CombinedLoggingHandler(os.Stdout, handler)

	for _, wrapper := range a.options.HandlerWrappers {
		h = wrapper(h)
	}

	a.mux.Handle(path, h)
}

func (a *httpApi) Run() error {
	if err := a.start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	log.Infof("http api received signal %s", <-ch)

	return a.stop()
}

func (a *httpApi) String() string {
	return "http"
}

func (a *httpApi) start() error {
	// TODO: tls
	listener, err := net.Listen("tcp", a.options.Address)
	if err != nil {
		return err
	}

	// make sure the address is right
	a.mtx.Lock()
	a.options.Address = listener.Addr().String()
	a.mtx.Unlock()

	log.Infof("http api is listening on %s", listener.Addr().String())

	go http.Serve(listener, a.mux)

	go func() {
		ch := <-a.exit
		ch <- listener.Close()
	}()

	return nil
}

func (a *httpApi) stop() error {
	ch := make(chan error)
	a.exit <- ch
	return <-ch
}

func NewApi(opts ...api.ApiOption) api.Api {
	options := api.NewApiOptions(opts...)

	a := &httpApi{
		options: options,
		mux:     http.NewServeMux(),
		mtx:     sync.RWMutex{},
		exit:    make(chan chan error),
	}

	return a
}
