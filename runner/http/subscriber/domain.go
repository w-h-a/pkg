package subscriber

import "github.com/w-h-a/pkg/sidecar"

type RouteEvent struct {
	Route string
	Event *sidecar.Event
}
