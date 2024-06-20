package authz

import (
	"context"

	"github.com/w-h-a/pkg/store"
)

type AuthzOption func(o *AuthzOptions)

type AuthzOptions struct {
	Store   store.Store
	Context context.Context
}

func AuthzWithStore(s store.Store) AuthzOption {
	return func(o *AuthzOptions) {
		o.Store = s
	}
}

func NewAuthzOptions(opts ...AuthzOption) AuthzOptions {
	options := AuthzOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
