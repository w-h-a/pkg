package mockclient

import (
	"context"
	"reflect"
	"sync"

	"github.com/w-h-a/pkg/client"
)

type mockStream struct {
	send Response
	recv Response
	err  error
	mtx  sync.RWMutex
}

func (s *mockStream) Context() context.Context {
	return nil
}

func (s *mockStream) Request() client.Request {
	return nil
}

func (s *mockStream) Send(_ interface{}) error {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	mock := s.send

	if mock.Err != nil {
		s.setError(mock.Err)
		return mock.Err
	}

	return nil
}

func (s *mockStream) Recv(msg interface{}) error {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	mock := s.recv

	if mock.Err != nil {
		s.setError(mock.Err)
		return mock.Err
	}

	val := reflect.ValueOf(msg)
	val = reflect.Indirect(val)

	response := mock.Response

	val.Set(reflect.ValueOf(response))

	return nil
}

func (s *mockStream) Error() error {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return s.err
}

func (s *mockStream) Close() error {
	return nil
}

func (s *mockStream) setError(e error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.err = e
}
