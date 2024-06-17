package client

import (
	"errors"

	"github.com/w-h-a/pkg/runtime"
)

var (
	ErrServiceNotFound = errors.New("service not found")
)

type Selector interface {
	Options() SelectorOptions
	Select(namespace, service string, port int, opts ...SelectOption) (func() (*runtime.Service, error), error)
	String() string
}
