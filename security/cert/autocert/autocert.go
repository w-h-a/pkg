package autocert

import (
	"net"

	"github.com/w-h-a/pkg/security/cert"
	goautocert "golang.org/x/crypto/acme/autocert"
)

type autocertProvider struct {
	options cert.CertOptions
}

func (c *autocertProvider) Options() cert.CertOptions {
	return c.options
}

func (c *autocertProvider) Listener(hosts ...string) (net.Listener, error) {
	return goautocert.NewListener(hosts...), nil
}

func (c *autocertProvider) String() string {
	return "autocert"
}

func NewCertProvider(opts ...cert.CertOption) cert.CertProvider {
	options := cert.NewCertOptions(opts...)

	a := &autocertProvider{
		options: options,
	}

	return a
}
