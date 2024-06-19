package authz

import (
	"context"

	"github.com/w-h-a/pkg/store"
)

type AuthzOption func(o *AuthzOptions)

type AuthzOptions struct {
	Store   store.Store
	Rules   []*Rule
	Context context.Context
}

func AuthzWithStore(s store.Store) AuthzOption {
	return func(o *AuthzOptions) {
		o.Store = s
	}
}

func AuthzWithRules(rs ...*Rule) AuthzOption {
	return func(o *AuthzOptions) {
		o.Rules = rs
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
