package secret

import "context"

type SecretOption func(o *SecretOptions)

type SecretOptions struct {
	Nodes   []string
	Prefix  string
	Context context.Context
}

func SecretWithNodes(addrs ...string) SecretOption {
	return func(o *SecretOptions) {
		o.Nodes = addrs
	}
}

func SecretWithPrefix(prefix string) SecretOption {
	return func(o *SecretOptions) {
		o.Prefix = prefix
	}
}

func NewSecretOptions(opts ...SecretOption) SecretOptions {
	options := SecretOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type GetSecretOption func(o *GetSecretOptions)

type GetSecretOptions struct {
	Context context.Context
}

func NewGetSecretOptions(opts ...GetSecretOption) GetSecretOptions {
	options := GetSecretOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
