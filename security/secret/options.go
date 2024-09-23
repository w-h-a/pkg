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
