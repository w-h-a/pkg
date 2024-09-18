package runner

type Process interface {
	Options() ProcessOptions
	Apply() error
	Destroy() error
	String() string
}
