package snssqs

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	transport "github.com/aws/smithy-go/endpoints"
	"github.com/w-h-a/pkg/broker"
)

type snsResolver struct {
	options broker.PublishOptions
}

func (r *snsResolver) ResolveEndpoint(ctx context.Context, params sns.EndpointParameters) (transport.Endpoint, error) {
	if len(r.options.Endpoint) == 0 {
		return sns.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	u, err := url.Parse(r.options.Endpoint)
	if err != nil {
		return transport.Endpoint{}, err
	}

	return transport.Endpoint{
		URI: *u,
	}, nil
}

type sqsResolver struct {
	options broker.SubscribeOptions
}

func (r *sqsResolver) ResolveEndpoint(ctx context.Context, params sqs.EndpointParameters) (transport.Endpoint, error) {
	if len(r.options.Endpoint) == 0 {
		return sqs.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	u, err := url.Parse(r.options.Endpoint)
	if err != nil {
		return transport.Endpoint{}, err
	}

	return transport.Endpoint{
		URI: *u,
	}, nil
}
