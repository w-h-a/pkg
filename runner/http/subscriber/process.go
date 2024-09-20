package subscriber

import (
	"context"
	"encoding/json"
	gohttp "net/http"
	"time"

	"github.com/w-h-a/pkg/runner"
	"github.com/w-h-a/pkg/runner/http"
	"github.com/w-h-a/pkg/sidecar"
)

type HttpSubscriber struct {
	proc  runner.Process
	event chan *RouteEvent
	exit  chan struct{}
}

func (p *HttpSubscriber) Options() runner.ProcessOptions {
	return p.proc.Options()
}

func (p *HttpSubscriber) Apply() error {
	return p.proc.Apply()
}

func (p *HttpSubscriber) Destroy() error {
	close(p.exit)
	return p.proc.Destroy()
}

func (p *HttpSubscriber) String() string {
	return "httpSubscriber"
}

func (p *HttpSubscriber) Receive() *RouteEvent {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil
	case routeEvent := <-p.event:
		return routeEvent
	}
}

func NewSubscriber(opts ...runner.ProcessOption) *HttpSubscriber {
	options := runner.NewProcessOptions(opts...)

	routes := []string{}

	if rs, ok := GetRoutesFromContext(options.Context); ok {
		routes = rs
	}

	event := make(chan *RouteEvent, 100)

	exit := make(chan struct{})

	for _, route := range routes {
		opts = append(opts, http.HttpProcessWithHandlerFuncs(route, func(w gohttp.ResponseWriter, r *gohttp.Request) {
			var sidecarEvent sidecar.Event

			if err := json.NewDecoder(r.Body).Decode(&sidecarEvent); err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}

			select {
			case <-exit:
				w.WriteHeader(500)
				return
			case <-r.Context().Done():
				w.WriteHeader(500)
				return
			case event <- &RouteEvent{Route: r.URL.Path, Event: &sidecarEvent}:
				w.WriteHeader(200)
				return
			}
		}))
	}

	opts = append(opts, http.HttpProcessWithHandlerFuncs("/health/check", func(w gohttp.ResponseWriter, r *gohttp.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}))

	s := &HttpSubscriber{
		proc:  http.NewProcess(opts...),
		event: event,
		exit:  exit,
	}

	return s
}
