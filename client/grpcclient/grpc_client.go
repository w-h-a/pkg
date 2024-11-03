package grpcclient

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/runtime"
	"github.com/w-h-a/pkg/telemetry/tracev2"
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

	callOptions := client.NewCallOptions(&c.options.CallOptions, opts...)

	next, err := c.next(req, callOptions)
	if err != nil {
		return err
	}

	// TODO: check if we already have a deadline
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

		name := req.Service()

		service, err := next()
		if err != nil {
			if err == client.ErrServiceNotFound {
				return errorutils.InternalServerError("client", "failed to find %s.%s: %v", name, namespace, err)
			}
			return errorutils.InternalServerError("client", "failed to select %s.%s: %v", name, namespace, err)
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

	// retry loop
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

func (c *grpcClient) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	callOptions := client.NewCallOptions(&c.options.CallOptions, opts...)

	next, err := c.next(req, callOptions)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, errorutils.Timeout("client", fmt.Sprintf("%v", ctx.Err()))
	default:
	}

	call := func(i int) (client.Stream, error) {
		duration, err := callOptions.Backoff(ctx, req, i)
		if err != nil {
			return nil, errorutils.InternalServerError("client", err.Error())
		}

		if duration.Seconds() > 0 {
			time.Sleep(duration)
		}

		namespace := req.Namespace()

		name := req.Service()

		service, err := next()
		if err != nil {
			if err == client.ErrServiceNotFound {
				return nil, errorutils.InternalServerError("client", "failed to find %s.%s: %v", name, namespace, err)
			}
			return nil, errorutils.InternalServerError("client", "failed to select %s.%s: %v", name, namespace, err)
		}

		// TODO: refactor this cruft
		address := service.Name + "." + service.Namespace + ":" + fmt.Sprintf("%d", service.Port)

		if len(service.Address) > 0 {
			address = service.Address
		}

		stream, err := c.stream(ctx, address, req, callOptions)
		if e, ok := err.(*errorutils.Error); ok {
			return stream, e
		}

		return stream, err
	}

	type response struct {
		stream client.Stream
		err    error
	}

	ch := make(chan response, callOptions.RetryCount+1)

	var e error

	// retry loop
	for i := 0; i <= callOptions.RetryCount; i++ {
		go func(i int) {
			s, err := call(i)
			ch <- response{s, err}
		}(i)

		select {
		case <-ctx.Done():
			return nil, errorutils.Timeout("client", fmt.Sprintf("%v", ctx.Err()))
		case rsp := <-ch:
			if rsp.err == nil {
				return rsp.stream, rsp.err
			}

			shouldRetry, retryErr := callOptions.RetryCheck(ctx, req, i, rsp.err)
			if retryErr != nil {
				return nil, retryErr
			}

			if !shouldRetry {
				return nil, rsp.err
			}

			e = rsp.err
		}
	}

	return nil, e
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

	if md, ok := metadatautils.FromContext(ctx); ok {
		for k, v := range md {
			header[k] = v
		}
	}

	if traceId, foundTrace := tracev2.TraceIdFromContext(ctx); foundTrace {
		if spanId, foundSpan := tracev2.SpanIdFromContext(ctx); foundSpan {
			header[tracev2.TraceParentKey] = fmt.Sprintf("00-%s-%s-01", hex.EncodeToString(traceId[:]), hex.EncodeToString(spanId[:]))
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

	conn, err := grpc.NewClient(address, grpcDialOptions...)
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

		err := conn.Invoke(
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

func (c *grpcClient) stream(ctx context.Context, address string, req client.Request, _ client.CallOptions) (client.Stream, error) {
	header := map[string]string{}

	md, ok := metadatautils.FromContext(ctx)
	if ok {
		for k, v := range md {
			header[k] = v
		}
	}

	header["content-type"] = req.ContentType()

	grpcMetadata := metadata.New(header)

	ctx = metadata.NewOutgoingContext(ctx, grpcMetadata)

	marshaler, err := c.newMarshaler(req.ContentType())
	if err != nil {
		return nil, errorutils.InternalServerError("client", err.Error())
	}

	grpcDialOptions := []grpc.DialOption{
		c.withCreds(address),
	}

	conn, err := grpc.NewClient(address, grpcDialOptions...)
	if err != nil {
		return nil, errorutils.InternalServerError("client", fmt.Sprintf("failed to get client connection: %v", err))
	}

	grpcCallOptions := []grpc.CallOption{
		grpc.ForceCodec(marshaler),
		grpc.CallContentSubtype(marshaler.Name()),
	}

	newCtx, cancel := context.WithCancel(ctx)

	s, err := conn.NewStream(
		newCtx,
		&grpc.StreamDesc{
			StreamName:    req.Method(),
			ClientStreams: true,
			ServerStreams: true,
		},
		ToGRPCMethod(req.Method()),
		grpcCallOptions...,
	)
	if err != nil {
		cancel()
		conn.Close()
		return nil, errorutils.InternalServerError("client", fmt.Sprintf("failed to create stream: %v", err))
	}

	return &grpcStream{
		context:    ctx,
		cancel:     cancel,
		request:    req,
		connection: conn,
		stream:     s,
		mtx:        sync.RWMutex{},
	}, nil
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
	}

	// need this for wrapping
	c := client.Client(g)
	for i := len(options.ClientWrappers); i > 0; i-- {
		c = options.ClientWrappers[i-1](c)
	}

	return c
}
