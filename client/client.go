package client

import (
	"context"
	"time"

	"github.com/w-h-a/pkg/utils/retryutils"
)

var (
	defaultBackoff = func(ctx context.Context, req Request, attempts int) (time.Duration, error) {
		return retryutils.ExponentialBackoff(attempts), nil
	}
	defaultRetryCheck = func(ctx context.Context, req Request, retryCount int, err error) (bool, error) {
		return retryutils.RetryOnError(err)
	}
	defaultRetryCount     = 1
	defaultRequestTimeout = time.Second * 5
)

type Client interface {
	Options() ClientOptions
	NewRequest(opts ...RequestOption) Request
	Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) (int, error)
	Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error)
	String() string
}
