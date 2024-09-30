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
	started bool
	exit    chan chan error
	mtx     sync.RWMutex
}

func (a *httpApi) Options() api.ApiOptions {
	return a.options
}

func (a *httpApi) Handle(path string, handler http.Handler) {
	// TODO: rm this
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
	a.mtx.RLock()
	if a.started {
		a.mtx.RUnlock()
		return nil
	}
	a.mtx.RUnlock()

	// TODO: log.SetLog(a.options.Logger)

	// TODO: register health handler

	// TODO: init exporters?

	var listener net.Listener

	var err error

	if a.options.EnableTLS && a.options.CertProvider != nil {
		// should we check the address to make sure it's :443?
		listener, err = a.options.CertProvider.Listener(a.options.Hosts...)
	} else {
		listener, err = net.Listen("tcp", a.options.Address)
	}

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

	a.mtx.Lock()
	a.started = true
	a.mtx.Unlock()

	return nil
}

func (a *httpApi) stop() error {
	a.mtx.RLock()
	if !a.started {
		a.mtx.RUnlock()
		return nil
	}
	a.mtx.RUnlock()

	ch := make(chan error)

	// signal start loop
	a.exit <- ch

	// wait for errors
	err := <-ch

	a.mtx.Lock()
	a.started = false
	a.mtx.Unlock()

	return err
}

func NewApi(opts ...api.ApiOption) api.Api {
	options := api.NewApiOptions(opts...)

	a := &httpApi{
		options: options,
		mux:     http.NewServeMux(),
		exit:    make(chan chan error),
		mtx:     sync.RWMutex{},
	}

	return a
}
