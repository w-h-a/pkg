package runner

type Manager interface {
	Options() ManagerOptions
	Apply() error
	Destroy() error
	String() string
}
