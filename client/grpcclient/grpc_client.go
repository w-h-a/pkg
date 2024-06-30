package grpcclient

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/w-h-a/pkg/broker"
	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/runtime"
	"github.com/w-h-a/pkg/utils/errorutils"
	"github.com/w-h-a/pkg/utils/marshalutils"
	"github.com/w-h-a/pkg/utils/metadatautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
)

var (
	defaultContentType = "application/grpc+proto"
)

type grpcClient struct {
	options client.ClientOptions
	once    *sync.Once
}

func (c *grpcClient) Options() client.ClientOptions {
	return c.options
}

func (c *grpcClient) NewRequest(opts ...client.RequestOption) client.Request {
	return NewRequest(opts...)
}

func (c *grpcClient) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if req == nil {
		return errorutils.InternalServerError("client", "req is nil")
	}

	if rsp == nil {
		return errorutils.InternalServerError("client", "rsp is nil")
	}

	callOptions := client.NewCallOptions(&c.options.CallOptions, opts...)

	next, err := c.next(req, callOptions)
	if err != nil {
		return err
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, callOptions.RequestTimeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return errorutils.Timeout("client", fmt.Sprintf("%v", ctx.Err()))
	default:
	}

	actualCall := c.call
	for i := len(callOptions.CallWrappers); i > 0; i-- {
		actualCall = callOptions.CallWrappers[i-1](actualCall)
	}

	call := func(i int) error {
		duration, err := callOptions.Backoff(ctx, req, i)
		if err != nil {
			return errorutils.InternalServerError("client", err.Error())
		}

		if duration.Seconds() > 0 {
			time.Sleep(duration)
		}

		namespace := req.Namespace()

		server := req.Service()

		service, err := next()
		if err != nil {
			if err == client.ErrServiceNotFound {
				return errorutils.InternalServerError("client", "failed to find %s.%s: %v", server, namespace, err)
			}
			return errorutils.InternalServerError("client", "failed to select %s.%s: %v", server, namespace, err)
		}

		// TODO: refactor this cruft
		address := service.Name + "." + service.Namespace + ":" + fmt.Sprintf("%d", service.Port)

		if len(service.Address) > 0 {
			address = service.Address
		}

		err = actualCall(ctx, address, req, rsp, callOptions)
		if e, ok := err.(*errorutils.Error); ok {
			return e
		}

		return err
	}

	ch := make(chan error, callOptions.RetryCount+1)

	var e error

	// retry lopp
	for i := 0; i <= callOptions.RetryCount; i++ {
		go func(i int) {
			ch <- call(i)
		}(i)

		select {
		case <-ctx.Done():
			return errorutils.Timeout("client", fmt.Sprintf("%v", ctx.Err()))
		case err := <-ch:
			if err == nil {
				return nil
			}

			shouldRetry, retryErr := callOptions.RetryCheck(ctx, req, i, err)
			if retryErr != nil {
				return retryErr
			}

			if !shouldRetry {
				return err
			}

			e = err
		}
	}

	return e
}

func (c *grpcClient) NewPublication(opts ...client.PublicationOption) client.Publication {
	return NewPublication(opts...)
}

func (c *grpcClient) Publish(ctx context.Context, pub client.Publication) error {
	header := map[string]string{}

	md, ok := metadatautils.FromContext(ctx)
	if ok {
		for k, v := range md {
			header[k] = v
		}
	}

	header["content-type"] = pub.ContentType()

	marshaler, err := c.newMarshaler(pub.ContentType())
	if err != nil {
		return errorutils.InternalServerError("client", err.Error())
	}

	bytes, err := marshaler.Marshal(pub.Unmarshaled())
	if err != nil {
		return errorutils.InternalServerError("client", err.Error())
	}

	c.once.Do(func() {
		c.options.Broker.Connect()
	})

	message := &broker.Message{
		Header: header,
		Body:   bytes,
	}

	return c.options.Broker.Publish(pub.Topic(), message)
}

func (c *grpcClient) String() string {
	return "grpc"
}

func (c *grpcClient) next(request client.Request, options client.CallOptions) (func() (*runtime.Service, error), error) {
	namespace := request.Namespace()
	name := request.Service()
	port := request.Port()

	// if we have the address already, use that
	if len(options.Address) > 0 {
		return func() (*runtime.Service, error) {
			return &runtime.Service{
				Namespace: namespace,
				Name:      name,
				Port:      port,
				Address:   options.Address,
			}, nil
		}, nil
	}

	// otherwise get the details from the selector
	next, err := c.options.Selector.Select(namespace, name, port, options.SelectOpts...)
	if err != nil {
		if err == client.ErrServiceNotFound {
			return nil, errorutils.InternalServerError("client", "failed to find %s.%s: %v", name, namespace, err)
		}
		return nil, errorutils.InternalServerError("client", "failed to select %s.%s: %v", name, namespace, err)
	}

	return next, nil
}

func (c *grpcClient) call(ctx context.Context, address string, req client.Request, rsp interface{}, options client.CallOptions) error {
	header := map[string]string{}

	md, ok := metadatautils.FromContext(ctx)
	if ok {
		for k, v := range md {
			header[k] = v
		}
	}

	header["timeout"] = fmt.Sprintf("%d", options.RequestTimeout)

	header["content-type"] = req.ContentType()

	delete(header, "connection")

	grpcMetadata := metadata.New(header)

	ctx = metadata.NewOutgoingContext(ctx, grpcMetadata)

	marshaler, err := c.newMarshaler(req.ContentType())
	if err != nil {
		return errorutils.InternalServerError("client", err.Error())
	}

	grpcDialOptions := []grpc.DialOption{
		c.withCreds(address),
	}

	clientConn, err := grpc.NewClient(address, grpcDialOptions...)
	if err != nil {
		return errorutils.InternalServerError("client", fmt.Sprintf("failed to get client connection: %v", err))
	}

	ch := make(chan error, 1)

	var e error

	go func() {
		grpcCallOptions := []grpc.CallOption{
			grpc.ForceCodec(marshaler),
			grpc.CallContentSubtype(marshaler.Name()),
		}

		err := clientConn.Invoke(
			ctx,
			ToGRPCMethod(req.Method()),
			req.Unmarshaled(),
			rsp,
			grpcCallOptions...,
		)

		ch <- err
	}()

	select {
	case err := <-ch:
		e = err
	case <-ctx.Done():
		e = errorutils.Timeout("client", "%v", ctx.Err())
	}

	return e
}

func (c *grpcClient) newMarshaler(contentType string) (encoding.Codec, error) {
	marshaler, ok := marshalutils.DefaultMarshalers[contentType]
	if !ok {
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	return marshaler, nil
}

func (c *grpcClient) withCreds(_ string) grpc.DialOption {
	// TODO
	return grpc.WithTransportCredentials(insecure.NewCredentials())
}

func init() {
	encoding.RegisterCodec(marshalutils.DefaultMarshalers["application/json"])
	encoding.RegisterCodec(marshalutils.DefaultMarshalers["application/proto"])
}

func NewClient(opts ...client.ClientOption) client.Client {
	options := client.NewClientOptions(opts...)

	if len(options.ContentType) == 0 {
		options.ContentType = defaultContentType
	}

	if options.Selector == nil {
		options.Selector = NewSelector()
	}

	g := &grpcClient{
		options: options,
		once:    &sync.Once{},
	}

	// need this for wrapping
	c := client.Client(g)
	for i := len(options.ClientWrappers); i > 0; i-- {
		c = options.ClientWrappers[i-1](c)
	}

	return c
}
