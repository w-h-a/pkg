package snssqs

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	transport "github.com/aws/smithy-go/endpoints"
)

type snsResolver struct {
	nodes []string
}

// TODO: figure out defaults
func (r *snsResolver) ResolveEndpoint(ctx context.Context, params sns.EndpointParameters) (transport.Endpoint, error) {
	u, err := url.Parse(r.nodes[0])
	if err != nil {
		return transport.Endpoint{}, err
	}

	return transport.Endpoint{
		URI: *u,
	}, nil
}

type sqsResolver struct {
	nodes []string
}

// TODO: figure out defaults
func (r *sqsResolver) ResolveEndpoint(ctx context.Context, params sqs.EndpointParameters) (transport.Endpoint, error) {
	u, err := url.Parse(r.nodes[0])
	if err != nil {
		return transport.Endpoint{}, err
	}

	return transport.Endpoint{
		URI: *u,
	}, nil
}
