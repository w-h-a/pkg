package handlers

import (
	"context"

	pb "github.com/w-h-a/pkg/proto/trace"
	"github.com/w-h-a/pkg/telemetry/trace"
	"github.com/w-h-a/pkg/utils/errorutils"
)

type TraceHandler interface {
	Read(ctx context.Context, req *pb.TraceRequest, rsp *pb.TraceResponse) error
}

type Trace struct {
	TraceHandler
}

type traceHandler struct {
	tracer trace.Trace
}

func (h *traceHandler) Read(ctx context.Context, req *pb.TraceRequest, rsp *pb.TraceResponse) error {
	spans, err := h.tracer.Read(
		trace.ReadWithTrace(req.Id),
		trace.ReadWithCount(int(req.Count)),
	)
	if err != nil {
		return errorutils.InternalServerError("trace", "failed to retrieve traces: %v", err)
	}

	for _, span := range spans {
		rsp.Spans = append(rsp.Spans, SerializeSpan(span))
	}

	return nil
}

func NewTraceHandler(tracer trace.Trace) TraceHandler {
	return &Trace{&traceHandler{
		tracer: tracer,
	}}
}
