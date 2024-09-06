package docker

import (
	"context"

	"github.com/w-h-a/pkg/runner"
)

type clientKey struct{}

func DockerRunnerWithClient(c DockerClient) runner.RunnerOption {
	return func(o *runner.RunnerOptions) {
		o.Context = context.WithValue(o.Context, clientKey{}, c)
	}
}

func GetClientFromContext(ctx context.Context) (DockerClient, bool) {
	c, ok := ctx.Value(clientKey{}).(DockerClient)
	return c, ok
}
