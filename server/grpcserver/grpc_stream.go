package grpcserver

import (
	"context"

	"github.com/w-h-a/pkg/server"
	"google.golang.org/grpc"
)

type grpcStream struct {
	request server.Request
	stream  grpc.ServerStream
}

func (s *grpcStream) Context() context.Context {
	return s.stream.Context()
}

func (s *grpcStream) Request() server.Request {
	return s.request
}

func (s *grpcStream) Recv(msg interface{}) error {
	return s.stream.RecvMsg(msg)
}

func (s *grpcStream) Send(msg interface{}) error {
	return s.stream.SendMsg(msg)
}
