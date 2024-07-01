package grpcclient

import (
	"context"
	"io"
	"sync"

	"github.com/w-h-a/pkg/client"
	"google.golang.org/grpc"
)

type grpcStream struct {
	context    context.Context
	cancel     func()
	request    client.Request
	connection *grpc.ClientConn
	stream     grpc.ClientStream
	err        error
	closed     bool
	mtx        sync.RWMutex
}

func (s *grpcStream) Context() context.Context {
	return s.context
}

func (s *grpcStream) Request() client.Request {
	return s.request
}

func (s *grpcStream) Send(msg interface{}) error {
	if err := s.stream.SendMsg(msg); err != nil {
		s.setError(err)
		return err
	}

	return nil
}

func (s *grpcStream) Recv(msg interface{}) error {
	if err := s.stream.RecvMsg(msg); err != nil {
		if err != io.EOF {
			s.setError(err)
		}
		return err
	}

	return nil
}

func (s *grpcStream) Error() error {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return s.err
}

func (s *grpcStream) Close() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true

	s.cancel()

	s.stream.CloseSend()

	return s.connection.Close()
}

func (s *grpcStream) setError(e error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.err = e
}
