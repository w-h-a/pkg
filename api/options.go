package api

import (
	"context"

	"github.com/w-h-a/pkg/security/cert"
)

type ApiOption func(o *ApiOptions)

type ApiOptions struct {
	Namespace       string
	Name            string
	Id              string
	Version         string
	Address         string
	Metadata        map[string]string
	EnableTLS       bool
	CertProvider    cert.CertProvider
	Hosts           []string
	HandlerWrappers []HandlerWrapper
	Context         context.Context
}

func ApiWithNamespace(ns string) ApiOption {
	return func(o *ApiOptions) {
		o.Namespace = ns
	}
}

func ApiWithName(n string) ApiOption {
	return func(o *ApiOptions) {
		o.Name = n
	}
}

func ApiWithId(id string) ApiOption {
	return func(o *ApiOptions) {
		o.Id = id
	}
}

func ApiWithVersion(v string) ApiOption {
	return func(o *ApiOptions) {
		o.Version = v
	}
}

func ApiWithAddress(addr string) ApiOption {
	return func(o *ApiOptions) {
		o.Address = addr
	}
}

func ApiWithMetadata(md map[string]string) ApiOption {
	return func(o *ApiOptions) {
		o.Metadata = md
	}
}

func ApiWithTLS() ApiOption {
	return func(o *ApiOptions) {
		o.EnableTLS = true
	}
}

func ApiWithCertProvider(p cert.CertProvider) ApiOption {
	return func(o *ApiOptions) {
		o.CertProvider = p
	}
}

func ApiWithHosts(hs ...string) ApiOption {
	return func(o *ApiOptions) {
		o.Hosts = hs
	}
}

func WrapHandler(w HandlerWrapper) ApiOption {
	return func(o *ApiOptions) {
		o.HandlerWrappers = append(o.HandlerWrappers, w)
	}
}

func NewApiOptions(opts ...ApiOption) ApiOptions {
	options := ApiOptions{
		Namespace: defaultNamespace,
		Name:      defaultName,
		Id:        defaultID,
		Version:   defaultVersion,
		Address:   defaultAddress,
		Metadata:  map[string]string{},
		Context:   context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
