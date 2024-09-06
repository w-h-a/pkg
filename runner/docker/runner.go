package docker

import (
	"sync"
	"testing"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/runner"
)

type dockerRunner struct {
	options     runner.RunnerOptions
	inactive    []runner.Manager
	inactiveMtx sync.RWMutex
	active      []runner.Manager
	activeMtx   sync.RWMutex
	client      DockerClient
}

func (r *dockerRunner) Options() runner.RunnerOptions {
	return r.options
}

func (r *dockerRunner) Start(m *testing.M) int {
	defer r.Stop()

	if err := r.register(); err != nil {
		return runner.FailExitCode
	}

	if err := r.apply(); err != nil {
		return runner.FailExitCode
	}

	return m.Run()
}

func (r *dockerRunner) Stop() {
	r.destroy()
}

func (r *dockerRunner) String() string {
	return "docker"
}

func (r *dockerRunner) register() error {
	for _, f := range r.options.Files {
		d := NewManager(
			runner.ManagerWithFile(f),
		)

		r.inactiveMtx.Lock()
		r.inactive = append(r.inactive, d)
		r.inactiveMtx.Unlock()
	}

	return nil
}

func (r *dockerRunner) apply() error {
	for d := r.dequeue(); d != nil; d = r.dequeue() {
		if err := d.Apply(); err != nil {
			log.Errorf("failed to apply %s: %v", d.Options().File.Path, err)
			return err
		}

		r.activeMtx.Lock()
		r.active = append(r.active, d)
		r.activeMtx.Unlock()

		log.Infof("successfully applied %s", d.Options().File.Path)
	}

	return nil
}

func (r *dockerRunner) destroy() {
	for d := r.pop(); d != nil; d = r.pop() {
		if err := d.Destroy(); err != nil {
			log.Errorf("failed to destroy %s: %v", d.Options().File.Path, err)
		} else {
			log.Infof("successfully destroyed %s", d.Options().File.Path)
		}
	}
}

func (r *dockerRunner) dequeue() runner.Manager {
	r.inactiveMtx.Lock()
	defer r.inactiveMtx.Unlock()

	if len(r.inactive) == 0 {
		return nil
	}

	d := r.inactive[0]

	r.inactive = r.inactive[1:]

	return d
}

func (r *dockerRunner) pop() runner.Manager {
	r.activeMtx.Lock()
	defer r.activeMtx.Unlock()

	if len(r.active) == 0 {
		return nil
	}

	d := r.active[len(r.active)-1]

	r.active = r.active[:len(r.active)-1]

	return d
}

func NewTestRunner(opts ...runner.RunnerOption) runner.TestRunner {
	options := runner.NewRunnerOptions(opts...)

	r := &dockerRunner{
		options:     options,
		inactive:    []runner.Manager{},
		inactiveMtx: sync.RWMutex{},
		active:      []runner.Manager{},
		activeMtx:   sync.RWMutex{},
	}

	if c, ok := GetClientFromContext(options.Context); ok {
		r.client = c
	} else {
		r.client = NewDockerClient()
	}

	return r
}
