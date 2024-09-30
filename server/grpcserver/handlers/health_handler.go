package handlers

import (
	"context"

	"github.com/w-h-a/pkg/proto/health"
	"github.com/w-h-a/pkg/server"
	"github.com/w-h-a/pkg/telemetry/log"
)

type HealthHandler interface {
	Check(ctx context.Context, req *health.HealthRequest, rsp *health.HealthResponse) error
	Log(ctx context.Context, req *health.LogRequest, rsp *health.LogResponse) error
}

type healthHandler struct {
	log log.Log
}

func (c *healthHandler) Check(ctx context.Context, req *health.HealthRequest, rsp *health.HealthResponse) error {
	rsp.Status = "ok"
	return nil
}

func (c *healthHandler) Log(ctx context.Context, req *health.LogRequest, rsp *health.LogResponse) error {
	opts := []log.ReadOption{}

	count := int(req.Count)
	if count > 0 {
		opts = append(opts, log.ReadWithCount(count))
	}

	records, err := c.log.Read(opts...)
	if err != nil {
		return err
	}

	protos := []*health.Record{}

	for _, record := range records {
		metadata := map[string]string{}
		for k, v := range record.Metadata {
			metadata[k] = v
		}
		proto := &health.Record{
			Timestamp: record.Timestamp.Unix(),
			Message:   record.Message.(string),
			Metadata:  metadata,
		}
		protos = append(protos, proto)
	}

	rsp.Records = protos

	return nil
}

func NewHealthHandler(name string) HealthHandler {
	return &healthHandler{
		log: log.GetLogger(),
	}
}

type Health struct {
	HealthHandler
}

func RegisterHealthHandler(s server.Server, handler HealthHandler, opts ...server.HandlerOption) error {
	return s.Handle(
		s.NewHandler(
			&Health{handler},
			opts...,
		),
	)
}
