package subscriber

import (
	"context"

	"github.com/w-h-a/pkg/runner"
)

type routesKey struct{}

func HttpSubscriberWithRoutes(routes ...string) runner.ProcessOption {
	return func(o *runner.ProcessOptions) {
		o.Context = context.WithValue(o.Context, routesKey{}, routes)
	}
}

func GetRoutesFromContext(ctx context.Context) ([]string, bool) {
	routes, ok := ctx.Value(routesKey{}).([]string)
	return routes, ok
}
