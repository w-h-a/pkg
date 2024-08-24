package runtime

type Runtime interface {
	Options() RuntimeOptions
	GetServices(namespace string, opts ...GetServicesOption) ([]*Service, error)
	IsServicePresent(name, namespace string) bool
	CreateService(name, namespace string, labels map[string]string) error
	UpdateDeployment(obj interface{})
	Start() error
	Stop() error
	String() string
}
