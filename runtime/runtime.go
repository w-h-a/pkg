package runtime

type Runtime interface {
	Options() RuntimeOptions
	GetServices(namespace string, opts ...GetServicesOption) ([]*Service, error)
	String() string
}
