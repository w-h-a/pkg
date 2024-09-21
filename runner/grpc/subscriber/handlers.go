package subscriber

import (
	"context"

	"github.com/w-h-a/pkg/proto/health"
)

type HealthHandler interface {
	Check(ctx context.Context, req *health.HealthRequest, rsp *health.HealthResponse) error
}

type healthHandler struct{}

func (c *healthHandler) Check(ctx context.Context, req *health.HealthRequest, rsp *health.HealthResponse) error {
	rsp.Status = "ok"
	return nil
}

func NewHealthHandler() HealthHandler {
	return &healthHandler{}
}

type Health struct {
	HealthHandler
}
