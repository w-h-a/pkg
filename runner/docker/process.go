package docker

import (
	"os/exec"
	"sync"

	"github.com/w-h-a/pkg/runner"
)

type dockerProcess struct {
	options runner.ProcessOptions
	upCmd   *exec.Cmd
	downCmd *exec.Cmd
	mtx     sync.RWMutex
}

func (p *dockerProcess) Options() runner.ProcessOptions {
	return p.options
}

func (p *dockerProcess) Apply() error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.upCmd = exec.Command(p.options.UpBinPath, p.options.UpArgs...)

	if err := runner.Outputs(p.upCmd); err != nil {
		return err
	}

	if err := p.upCmd.Run(); err != nil {
		return err
	}

	return nil
}

func (p *dockerProcess) Destroy() error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.downCmd = exec.Command(p.options.DownBinPath, p.options.DownArgs...)

	if err := runner.Outputs(p.downCmd); err != nil {
		return err
	}

	if err := p.downCmd.Run(); err != nil {
		return err
	}

	return nil
}

func (p *dockerProcess) String() string {
	return "docker"
}

func NewProcess(opts ...runner.ProcessOption) runner.Process {
	options := runner.NewProcessOptions(opts...)

	p := &dockerProcess{
		options: options,
		mtx:     sync.RWMutex{},
	}

	return p
}
