package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/runtime"
	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/utils/errorutils"
	"github.com/w-h-a/pkg/utils/marshalutils"
	"github.com/w-h-a/pkg/utils/metadatautils"
)

const (
	defaultContentType = "application/json"
)

type httpClient struct {
	options client.ClientOptions
}

func (c *httpClient) Options() client.ClientOptions {
	return c.options
}

func (c *httpClient) NewRequest(opts ...client.RequestOption) client.Request {
	return NewRequest(opts...)
}

func (c *httpClient) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
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

func (c *httpClient) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	return nil, nil
}

func (c *httpClient) String() string {
	return "http"
}

func (c *httpClient) next(request client.Request, options client.CallOptions) (func() (*runtime.Service, error), error) {
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

func (c *httpClient) call(ctx context.Context, address string, req client.Request, rsp interface{}, options client.CallOptions) error {
	header := http.Header{}

	if md, ok := metadatautils.FromContext(ctx); ok {
		for k, v := range md {
			header.Set(k, v)
		}
	}

	header.Set("timeout", fmt.Sprintf("%d", options.RequestTimeout))

	header.Set("content-type", req.ContentType())

	marshaler, err := c.newMarshaler(req.ContentType())
	if err != nil {
		return errorutils.InternalServerError("client", err.Error())
	}

	bs, err := marshaler.Marshal(req.Unmarshaled())
	if err != nil {
		return errorutils.InternalServerError("client", err.Error())
	}

	buf := &buffer{bytes.NewBuffer(bs)}
	defer buf.Close()

	endpoint := req.Method()

	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	rawurl := "http://" + address + endpoint

	URL, err := url.Parse(rawurl)
	if err != nil {
		return errorutils.InternalServerError("client", err.Error())
	}

	httpReq := &http.Request{
		Method:        "POST",
		URL:           URL,
		Header:        header,
		Body:          buf,
		ContentLength: int64(len(bs)),
		Host:          address,
	}

	httpRsp, err := http.DefaultClient.Do(httpReq.WithContext(ctx))
	if err != nil {
		return errorutils.InternalServerError("client", err.Error())
	}

	defer httpRsp.Body.Close()

	bs, err = io.ReadAll(httpRsp.Body)
	if err != nil {
		return errorutils.InternalServerError("client", err.Error())
	}

	log.Infof("RECEIVED %d bytes", len(bs))

	log.Infof("RECEIVED %s", string(bs))

	log.Infof("RSP BEFORE %+v", rsp)

	if err := marshaler.Unmarshal(bs, rsp); err != nil {
		return errorutils.InternalServerError("client", err.Error())
	}

	log.Infof("RSP AFTER %+v", rsp)

	return nil
}

func (c *httpClient) newMarshaler(contentType string) (marshalutils.Marshaler, error) {
	marshaler, ok := marshalutils.DefaultMarshalers[contentType]
	if !ok {
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	return marshaler, nil
}

func NewClient(opts ...client.ClientOption) client.Client {
	options := client.NewClientOptions(opts...)

	if len(options.ContentType) == 0 {
		options.ContentType = defaultContentType
	}

	if options.Selector == nil {
		options.Selector = NewSelector()
	}

	h := &httpClient{
		options: options,
	}

	// wrap in reverse
	c := client.Client(h)
	for i := len(options.ClientWrappers); i > 0; i-- {
		c = options.ClientWrappers[i-1](c)
	}

	return c
}
