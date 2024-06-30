package main

import (
	"context"

	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/client/grpcclient"
	"github.com/w-h-a/pkg/proto/health"
	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/utils/metadatautils"
)


type clientWrapper struct {
	client.Client
	headers metadatautils.Metadata
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	newCtx := metadatautils.MergeContext(ctx, c.headers, false)

	return c.Client.Call(newCtx, req, rsp, opts...)
}

func NewClientWrapper(md metadatautils.Metadata) client.ClientWrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{c, md}
	}
}

type logWrapper struct {
	client.Client
}

func (c *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	log.Infof("making a call to %s of %s.%s", req.Method(), req.Service(), req.Namespace())
	
	return c.Client.Call(ctx, req, rsp, opts...)
}

func NewLogWrapper(c client.Client) client.Client {
	return &logWrapper{c}
}

func main() {
	call()
}

func call() {
	h := map[string]string{}

	grpcClient := grpcclient.NewClient(
		client.WrapClient(NewClientWrapper(h)),
		client.WrapClient(NewLogWrapper),
	)

	req := grpcClient.NewRequest(
		client.RequestWithNamespace("app"),
		client.RequestWithName("greeter"),
		client.RequestWithMethod("Health.Log"),
		client.RequestWithUnmarshaledRequest(
			&health.LogRequest{
				Count: int64(10),
			},
		),
	)

	rsp := &health.LogResponse{}

	if err := grpcClient.Call(context.Background(), req, rsp, client.CallWithAddress("127.0.0.1:53116")); err != nil {
		log.Fatalf("failed to make call: %v", err)
	}

	log.Infof("SUCCESS: %v", rsp.Records)
}