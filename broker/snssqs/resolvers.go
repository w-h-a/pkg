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

func (r *snsResolver) ResolveEndpoint(ctx context.Context, params sns.EndpointParameters) (transport.Endpoint, error) {
	if len(r.nodes) == 0 {
		return sns.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

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

func (r *sqsResolver) ResolveEndpoint(ctx context.Context, params sqs.EndpointParameters) (transport.Endpoint, error) {
	if len(r.nodes) == 0 {
		return sqs.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	u, err := url.Parse(r.nodes[0])
	if err != nil {
		return transport.Endpoint{}, err
	}

	return transport.Endpoint{
		URI: *u,
	}, nil
}
