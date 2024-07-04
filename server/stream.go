package server

import "context"

type Stream interface {
	Context() context.Context
	Request() Request
	Recv(msg interface{}) error
	Send(msg interface{}) error
	Error() error
	Close() error
}
