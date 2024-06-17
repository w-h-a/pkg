package controllers

import (
	"context"

	"github.com/w-h-a/pkg/proto/health"
	"github.com/w-h-a/pkg/server"
	"github.com/w-h-a/pkg/telemetry/log"
)

type HealthController interface {
	Check(ctx context.Context, req *health.HealthRequest, rsp *health.HealthResponse) error
	Log(ctx context.Context, req *health.LogRequest, rsp *health.LogResponse) error
}

type healthController struct {
	log log.Log
}

func (c *healthController) Check(ctx context.Context, req *health.HealthRequest, rsp *health.HealthResponse) error {
	log.Info("RESPONDING TO HEALTH CHECK")
	rsp.Status = "ok"
	return nil
}

func (c *healthController) Log(ctx context.Context, req *health.LogRequest, rsp *health.LogResponse) error {
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

func NewHealthController(name string) HealthController {
	log.SetName(name)

	return &healthController{
		log: log.GetLogger(),
	}
}

type Health struct {
	HealthController
}

func RegisterHealthController(s server.Server, controller HealthController, opts ...server.ControllerOption) error {
	return s.RegisterController(
		s.NewController(
			&Health{controller},
			opts...,
		),
	)
}
