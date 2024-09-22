package ssm

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	transport "github.com/aws/smithy-go/endpoints"
)

type ssmResolver struct {
	nodes []string
}

// TODO: figure out defaults
func (r *ssmResolver) ResolveEndpoint(ctx context.Context, params ssm.EndpointParameters) (transport.Endpoint, error) {
	u, err := url.Parse(r.nodes[0])
	if err != nil {
		return transport.Endpoint{}, err
	}

	return transport.Endpoint{
		URI: *u,
	}, nil
}
