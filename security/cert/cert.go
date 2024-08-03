package cert

import "net"

type CertProvider interface {
	Options() CertOptions
	Listener(hosts ...string) (net.Listener, error)
	String() string
}
