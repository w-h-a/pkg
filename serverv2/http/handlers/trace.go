package handlers

import (
	"net/http"
	"strconv"

	"github.com/w-h-a/pkg/telemetry/trace"
	"github.com/w-h-a/pkg/utils/errorutils"
	"github.com/w-h-a/pkg/utils/httputils"
)

type TraceHandler interface {
	Read(w http.ResponseWriter, r *http.Request)
}

type traceHandler struct {
	tracer trace.Trace
}

func (h *traceHandler) Read(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	count := r.URL.Query().Get("count")

	c, err := strconv.Atoi(count)
	if err != nil {
		httputils.ErrResponse(w, errorutils.BadRequest("trace", "received bad count query param: %v", err))
		return
	}

	spans, err := h.tracer.Read(
		trace.ReadWithTrace(id),
		trace.ReadWithCount(c),
	)
	if err != nil {
		httputils.ErrResponse(w, errorutils.InternalServerError("trace", "failed to retrieve traces: %v", err))
		return
	}

	httputils.OkResponse(w, spans)
}

func NewTraceHandler(tracer trace.Trace) TraceHandler {
	return &traceHandler{tracer}
}
