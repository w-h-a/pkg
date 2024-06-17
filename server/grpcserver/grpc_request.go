package grpcserver

import (
	"github.com/w-h-a/pkg/server"
)

type grpcRequest struct {
	options server.RequestOptions
}

func (r *grpcRequest) Options() server.RequestOptions {
	return r.options
}

func (r *grpcRequest) Namespace() string {
	return r.options.Namespace
}

func (r *grpcRequest) Server() string {
	return r.options.Name
}

func (r *grpcRequest) Method() string {
	return r.options.Method
}

func (r *grpcRequest) ContentType() string {
	return r.options.ContentType
}

func (r *grpcRequest) Unmarshaled() interface{} {
	return r.options.UnmarshaledRequest
}

func (r *grpcRequest) Marshaled() []byte {
	return r.options.MarshaledRequest
}

func (r *grpcRequest) String() string {
	return "grpc"
}

func NewRequest(opts ...server.RequestOption) server.Request {
	options := server.NewRequestOptions(opts...)

	return &grpcRequest{
		options: options,
	}
}
