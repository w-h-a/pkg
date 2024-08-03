package cert

import "net"

type CertProvider interface {
	Options() CertOptions
	Listener(domains ...string) (net.Listener, error)
	String() string
}
