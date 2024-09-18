package runner

import "testing"

const (
	SuccessExitCode = 0
	FailExitCode    = 1
)

type TestRunner interface {
	Options() RunnerOptions
	Start(m *testing.M) int
	Stop()
	String() string
}
