package grpcclient

import "github.com/w-h-a/pkg/client"

type grpcPublication struct {
	options client.PublicationOptions
}

func (p *grpcPublication) Options() client.PublicationOptions {
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

func NewPublication(opts ...client.PublicationOption) client.Publication {
	options := client.NewPublicationOptions(opts...)

	return &grpcPublication{
		options: options,
	}
}
