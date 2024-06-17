package grpcclient

import "github.com/w-h-a/pkg/client"

type grpcRequest struct {
	options client.RequestOptions
}

func (r *grpcRequest) Options() client.RequestOptions {
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

func (r *grpcRequest) Port() int {
	return r.options.Port
}

func (r *grpcRequest) ContentType() string {
	return r.options.ContentType
}

func (r *grpcRequest) Unmarshaled() interface{} {
	return r.options.UnmarshaledRequest
}

func (r *grpcRequest) String() string {
	return "grpc"
}

func NewRequest(opts ...client.RequestOption) client.Request {
	options := client.NewRequestOptions(opts...)

	if len(options.ContentType) == 0 {
		options.ContentType = defaultContentType
	}

	return &grpcRequest{
		options: options,
	}
}
