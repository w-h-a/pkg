package runtime

import (
	"context"
	"net/http"
)

type RuntimeOption func(o *RuntimeOptions)

type RuntimeOptions struct {
	Host        string
	BearerToken string
	Client      *http.Client
	Context     context.Context
}

func RuntimeWithHost(h string) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.Host = h
	}
}

func RuntimeWithBearerToken(t string) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.BearerToken = t
	}
}

func RuntimeWithClient(c *http.Client) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.Client = c
	}
}

func NewRuntimeOptions(opts ...RuntimeOption) RuntimeOptions {
	options := RuntimeOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type GetServicesOption func(o *GetServicesOptions)

type GetServicesOptions struct {
	Name    string
	Version string
}

func GetServicesWithName(n string) GetServicesOption {
	return func(o *GetServicesOptions) {
		o.Name = n
	}
}

func GetServicesWithVersion(v string) GetServicesOption {
	return func(o *GetServicesOptions) {
		o.Version = v
	}
}

func NewGetServicesOptions(opts ...GetServicesOption) GetServicesOptions {
	options := GetServicesOptions{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
