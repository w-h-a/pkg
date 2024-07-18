package secret

import "context"

type SecretOption func(o *SecretOptions)

type SecretOptions struct {
	Context context.Context
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

type EncryptOption func(o *EncryptOptions)

type EncryptOptions struct {
	Context context.Context
}

func NewEncryptOptions(opts ...EncryptOption) EncryptOptions {
	options := EncryptOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type DecryptOption func(o *DecryptOptions)

type DecryptOptions struct {
	Context context.Context
}

func NewDecryptOptions(opts ...DecryptOption) DecryptOptions {
	options := DecryptOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
