package authn

import (
	"context"
	"time"

	"github.com/w-h-a/pkg/store"
)

type AuthnOption func(o *AuthnOptions)

type AuthnOptions struct {
	Store   store.Store
	Context context.Context
}

func AuthnWithStore(s store.Store) AuthnOption {
	return func(o *AuthnOptions) {
		o.Store = s
	}
}

func NewAuthnOptions(opts ...AuthnOption) AuthnOptions {
	options := AuthnOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type GenerateOption func(o *GenerateOptions)

type GenerateOptions struct {
	Secret   string
	Roles    []string
	Metadata map[string]string
}

func GenerateWithSecret(s string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Secret = s
	}
}

func GenerateWithRoles(rs ...string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Roles = rs
	}
}

func GenerateWithMetadata(md map[string]string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Metadata = md
	}
}

func NewGenerateOptions(opts ...GenerateOption) GenerateOptions {
	options := GenerateOptions{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type TokenOption func(o *TokenOptions)

// TODO: optionally pass refresh token to get new access token instead of credentials
type TokenOptions struct {
	Id     string
	Secret string
	Expiry time.Duration
}

func TokenWithCredentials(id, secret string) TokenOption {
	return func(o *TokenOptions) {
		o.Id = id
		o.Secret = secret
	}
}

func TokenWithExpiry(expiry time.Duration) TokenOption {
	return func(o *TokenOptions) {
		o.Expiry = expiry
	}
}

func NewTokenOptions(opts ...TokenOption) TokenOptions {
	options := TokenOptions{
		Expiry: time.Minute,
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
