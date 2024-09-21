package subscriber

import (
	"context"
	"time"

	pb "github.com/w-h-a/pkg/proto/sidecar"
	"github.com/w-h-a/pkg/runner"
	"github.com/w-h-a/pkg/runner/grpc"
	grpcserver "github.com/w-h-a/pkg/serverv2/grpc"
)

type GrpcSubscriber struct {
	proc  runner.Process
	event chan *pb.Event
}

func (p *GrpcSubscriber) Options() runner.ProcessOptions {
	return p.proc.Options()
}

func (p *GrpcSubscriber) Apply() error {
	return p.proc.Apply()
}

func (p *GrpcSubscriber) Destroy() error {
	close(p.event)
	return p.proc.Destroy()
}

func (p *GrpcSubscriber) String() string {
	return "GrpcSubscriber"
}

func (p *GrpcSubscriber) Receive() *pb.Event {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil
	case event := <-p.event:
		return event
	}
}

func NewSubscriber(opts ...runner.ProcessOption) *GrpcSubscriber {
	event := make(chan *pb.Event, 100)

	opts = append(
		opts,
		grpc.GrpcProcessWithHandlers(
			grpcserver.NewHandler(
				&Health{NewHealthHandler()},
			),
		),
	)

	s := &GrpcSubscriber{
		proc:  grpc.NewProcess(opts...),
		event: event,
	}

	return s
}
