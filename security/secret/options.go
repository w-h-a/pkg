package secret

import "context"

type SecretOption func(o *SecretOptions)

type SecretOptions struct {
	Nodes   []string
	Context context.Context
}

func SecretWithNodes(addrs ...string) SecretOption {
	return func(o *SecretOptions) {
		o.Nodes = addrs
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
	Prefix  string
	Context context.Context
}

func GetSecretWithPrefix(prefix string) GetSecretOption {
	return func(o *GetSecretOptions) {
		o.Prefix = prefix
	}
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
