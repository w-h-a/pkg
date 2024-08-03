package cert

import "context"

type CertOption func(o *CertOptions)

type CertOptions struct {
	Context context.Context
}

func NewCertOptions(opts ...CertOption) CertOptions {
	options := CertOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
