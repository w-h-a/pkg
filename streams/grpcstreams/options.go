package grpcstreams

import (
	"context"

	"github.com/w-h-a/pkg/security/token"
	"github.com/w-h-a/pkg/store"
)

type storeKey struct{}

func GrpcStreamsWithStore(s store.Store) token.TokenOption {
	return func(o *token.TokenOptions) {
		o.Context = context.WithValue(o.Context, storeKey{}, s)
	}
}

func GetStoreFromContext(ctx context.Context) (store.Store, bool) {
	s, ok := ctx.Value(storeKey{}).(store.Store)
	return s, ok
}
