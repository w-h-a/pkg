package runner

import (
	"testing"
)

const (
	FailExitCode = 1
)

type TestRunner interface {
	Options() RunnerOptions
	Start(m *testing.M) int
	Stop()
	String() string
}
