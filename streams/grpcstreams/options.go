package grpcstreams

import (
	"context"

	"github.com/w-h-a/pkg/store"
	"github.com/w-h-a/pkg/streams"
)

type storeKey struct{}

func GrpcStreamsWithStore(s store.Store) streams.StreamsOption {
	return func(o *streams.StreamsOptions) {
		o.Context = context.WithValue(o.Context, storeKey{}, s)
	}
}

func GetStoreFromContext(ctx context.Context) (store.Store, bool) {
	s, ok := ctx.Value(storeKey{}).(store.Store)
	return s, ok
}
