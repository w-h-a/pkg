package runner

import (
	"sync"
	"testing"

	"github.com/w-h-a/pkg/telemetry/log"
)

const (
	SuccessExitCode = 0
	FailExitCode    = 1
)

type TestRunner struct {
	options     RunnerOptions
	inactive    []Process
	inactiveMtx sync.RWMutex
	active      []Process
	activeMtx   sync.RWMutex
}

func (r *TestRunner) Options() RunnerOptions {
	return r.options
}

func (r *TestRunner) Start(m *testing.M) int {
	defer r.Stop()

	if err := r.Register(); err != nil {
		return FailExitCode
	}

	if err := r.Apply(); err != nil {
		return FailExitCode
	}

	return m.Run()
}

func (r *TestRunner) Stop() {
	r.Destroy()
}

func (r *TestRunner) Register() error {
	for _, p := range r.options.Processes {
		r.inactiveMtx.Lock()
		r.inactive = append(r.inactive, p)
		r.inactiveMtx.Unlock()
	}

	return nil
}

func (r *TestRunner) Apply() error {
	for p := r.dequeue(); p != nil; p = r.dequeue() {
		if err := p.Apply(); err != nil {
			log.Errorf("failed to apply %s: %v", p.Options().Id, err)
			return err
		}

		r.activeMtx.Lock()
		r.active = append(r.active, p)
		r.activeMtx.Unlock()

		log.Infof("successfully applied %s", p.Options().Id)
	}

	return nil
}

func (r *TestRunner) Destroy() {
	for p := r.pop(); p != nil; p = r.pop() {
		if err := p.Destroy(); err != nil {
			log.Errorf("failed to destroy process %s: %v", p.Options().Id, err)
		} else {
			log.Infof("successfully destroyed process %s", p.Options().Id)
		}
	}
}

func (r *TestRunner) dequeue() Process {
	r.inactiveMtx.Lock()
	defer r.inactiveMtx.Unlock()

	if len(r.inactive) == 0 {
		return nil
	}

	m := r.inactive[0]

	r.inactive = r.inactive[1:]

	return m
}

func (r *TestRunner) pop() Process {
	r.activeMtx.Lock()
	defer r.activeMtx.Unlock()

	if len(r.active) == 0 {
		return nil
	}

	m := r.active[len(r.active)-1]

	r.active = r.active[:len(r.active)-1]

	return m
}

func NewTestRunner(opts ...RunnerOption) *TestRunner {
	options := NewRunnerOptions(opts...)

	r := &TestRunner{
		options:     options,
		inactive:    []Process{},
		inactiveMtx: sync.RWMutex{},
		active:      []Process{},
		activeMtx:   sync.RWMutex{},
	}

	return r
}
