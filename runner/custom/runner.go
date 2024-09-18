package custom

import (
	"sync"
	"testing"

	"github.com/w-h-a/pkg/runner"
	"github.com/w-h-a/pkg/telemetry/log"
)

type customRunner struct {
	options     runner.RunnerOptions
	inactive    []runner.Process
	inactiveMtx sync.RWMutex
	active      []runner.Process
	activeMtx   sync.RWMutex
}

func (r *customRunner) Options() runner.RunnerOptions {
	return r.options
}

func (r *customRunner) Start(m *testing.M) int {
	defer r.Stop()

	if err := r.register(); err != nil {
		return runner.FailExitCode
	}

	if err := r.apply(); err != nil {
		return runner.FailExitCode
	}

	return m.Run()
}

func (r *customRunner) Stop() {
	r.destroy()
}

func (r *customRunner) String() string {
	return "custom"
}

func (r *customRunner) register() error {
	for _, p := range r.options.Processes {
		r.inactiveMtx.Lock()
		r.inactive = append(r.inactive, p)
		r.inactiveMtx.Unlock()
	}

	return nil
}

func (r *customRunner) apply() error {
	for p := r.dequeue(); p != nil; p = r.dequeue() {
		if err := p.Apply(); err != nil {
			log.Errorf("failed to apply %s: %v", p.Options().UpBinPath, err)
			return err
		}

		r.activeMtx.Lock()
		r.active = append(r.active, p)
		r.activeMtx.Unlock()

		log.Infof("successfully applied %s", p.Options().UpBinPath)
	}

	return nil
}

func (r *customRunner) destroy() {
	for p := r.pop(); p != nil; p = r.pop() {
		if err := p.Destroy(); err != nil {
			log.Errorf("failed to destroy process: %v", err)
		} else {
			log.Infof("successfully destroyed process")
		}
	}
}

func (r *customRunner) dequeue() runner.Process {
	r.inactiveMtx.Lock()
	defer r.inactiveMtx.Unlock()

	if len(r.inactive) == 0 {
		return nil
	}

	m := r.inactive[0]

	r.inactive = r.inactive[1:]

	return m
}

func (r *customRunner) pop() runner.Process {
	r.activeMtx.Lock()
	defer r.activeMtx.Unlock()

	if len(r.active) == 0 {
		return nil
	}

	m := r.active[len(r.active)-1]

	r.active = r.active[:len(r.active)-1]

	return m
}

func NewTestRunner(opts ...runner.RunnerOption) runner.TestRunner {
	options := runner.NewRunnerOptions(opts...)

	r := &customRunner{
		options:     options,
		inactive:    []runner.Process{},
		inactiveMtx: sync.RWMutex{},
		active:      []runner.Process{},
		activeMtx:   sync.RWMutex{},
	}

	return r
}
