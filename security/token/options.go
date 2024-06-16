package token

import (
	"context"
	"time"
)

type TokenOption func(o *TokenOptions)

type TokenOptions struct {
	Context context.Context
}

func NewTokenOptions(opts ...TokenOption) TokenOptions {
	options := TokenOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type GenerateOption func(o *GenerateOptions)

type GenerateOptions struct {
	Expiry   time.Duration
	Id       string
	Roles    []string
	Metadata map[string]string
}

func GenerateWithExpiry(d time.Duration) GenerateOption {
	return func(o *GenerateOptions) {
		o.Expiry = d
	}
}

func GenerateWithId(id string) GenerateOption {
	return func(o *GenerateOptions) {
		o.Id = id
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
	options := GenerateOptions{
		Expiry: time.Minute * 15,
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
