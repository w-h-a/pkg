package grpcserver

import "github.com/w-h-a/pkg/server"

type grpcPublication struct {
	options server.PublicationOptions
}

func (p *grpcPublication) Options() server.PublicationOptions {
	return p.options
}

func (p *grpcPublication) Topic() string {
	return p.options.Topic
}

func (p *grpcPublication) ContentType() string {
	return p.options.ContentType
}

func (p *grpcPublication) Unmarshaled() interface{} {
	return p.options.UnmarshaledPayload
}

func (p *grpcPublication) String() string {
	return "grpc"
}

func NewPublication(opts ...server.PublicationOption) server.Publication {
	options := server.NewPublicationOptions(opts...)

	return &grpcPublication{
		options: options,
	}
}
