package client

import "context"

type Stream interface {
	Context() context.Context
	Request() Request
	Send(msg interface{}) error
	Recv(msg interface{}) error
	Error() error
	Close() error
}
