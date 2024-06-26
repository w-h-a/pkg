package mockclient

import (
	"context"

	"github.com/w-h-a/pkg/client"
)

type responsesKey struct{}
type clientKey struct{}

func ClientWithResponses(service, method string, response Response) client.ClientOption {
	return func(o *client.ClientOptions) {
		responses, ok := GetResponsesFromContext(o.Context)
		if !ok {
			responses = map[string]Response{}
		}

		responses[service+":"+method] = response

		o.Context = context.WithValue(o.Context, responsesKey{}, responses)
	}
}

func GetResponsesFromContext(ctx context.Context) (map[string]Response, bool) {
	rsp, ok := ctx.Value(responsesKey{}).(map[string]Response)
	return rsp, ok
}

func ClientWithClient(c client.Client) client.ClientOption {
	return func(o *client.ClientOptions) {
		o.Context = context.WithValue(o.Context, clientKey{}, c)
	}
}

func GetClientFromContext(ctx context.Context) (client.Client, bool) {
	c, ok := ctx.Value(clientKey{}).(client.Client)
	return c, ok
}
