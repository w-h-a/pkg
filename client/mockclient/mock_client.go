package mockclient

import (
	"context"
	"reflect"
	"sync"

	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/client/grpcclient"
	"github.com/w-h-a/pkg/utils/errorutils"
)

type mockClient struct {
	options   client.ClientOptions
	responses map[string]Response
	streams   map[string]client.Stream
	client    client.Client
	mtx       sync.RWMutex
}

func (c *mockClient) Options() client.ClientOptions {
	return c.options
}

func (c *mockClient) NewRequest(opts ...client.RequestOption) client.Request {
	return c.client.NewRequest(opts...)
}

func (c *mockClient) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	mock, ok := c.responses[req.Service()+":"+req.Method()]
	if !ok {
		return errorutils.NotFound("mock.client", "service:method %s:%s not found in responses %+v", req.Service(), req.Method(), c.responses)
	}

	if mock.Err != nil {
		return mock.Err
	}

	val := reflect.ValueOf(rsp)
	val = reflect.Indirect(val)

	response := mock.Response

	val.Set(reflect.ValueOf(response))

	return nil
}

func (c *mockClient) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	mock, ok := c.streams[req.Service()+":"+req.Method()]
	if !ok {
		return nil, errorutils.NotFound("mock.client", "service:method %s:%s not found in streams %+v", req.Service(), req.Method(), c.streams)
	}

	return mock, nil
}

func (c *mockClient) String() string {
	return "mock"
}

func NewClient(opts ...client.ClientOption) client.Client {
	options := client.NewClientOptions(opts...)

	responses, ok := GetResponsesFromContext(options.Context)
	if !ok {
		responses = map[string]Response{}
	}

	streams, ok := GetStreamsFromContext(options.Context)
	if !ok {
		streams = map[string]client.Stream{}
	}

	c, ok := GetClientFromContext(options.Context)
	if !ok {
		c = grpcclient.NewClient()
	}

	m := &mockClient{
		options:   options,
		responses: responses,
		streams:   streams,
		client:    c,
		mtx:       sync.RWMutex{},
	}

	return m
}
