package handlers

import (
	pb "github.com/w-h-a/pkg/proto/trace"
	"github.com/w-h-a/pkg/telemetry/trace"
)

func SerializeSpan(s *trace.Span) *pb.Span {
	return &pb.Span{
		Name:     s.Name,
		Id:       s.Id,
		Parent:   s.Parent,
		Trace:    s.Trace,
		Started:  uint64(s.Started.UnixNano()),
		Duration: uint64(s.Duration.Nanoseconds()),
		Metadata: s.Metadata,
	}
}
