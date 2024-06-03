package runtime

type Runtime interface {
	Options() RuntimeOptions
	GetServices(opts ...GetServicesOption) ([]*Service, error)
	String() string
}
