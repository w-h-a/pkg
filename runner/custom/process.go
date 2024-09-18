package custom

import (
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/w-h-a/pkg/runner"
)

type customProcess struct {
	options runner.ProcessOptions
	upCmd   *exec.Cmd
	downCmd *exec.Cmd
	mtx     sync.RWMutex
}

func (p *customProcess) Options() runner.ProcessOptions {
	return p.options
}

func (p *customProcess) Apply() error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.upCmd = exec.Command(p.options.UpBinPath, p.options.UpArgs...)

	if err := runner.Outputs(p.upCmd); err != nil {
		return err
	}

	for k, v := range p.options.EnvVars {
		p.upCmd.Env = append(p.upCmd.Env, k+"="+v)
	}

	errCh := make(chan error)

	go func() {
		err := p.upCmd.Start()
		for p.upCmd.Process == nil {
			time.Sleep(1 * time.Second)
		}
		errCh <- err
	}()

	return <-errCh
}

func (p *customProcess) Destroy() error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.upCmd != nil && p.upCmd.Process != nil {
		if err := p.upCmd.Process.Signal(os.Interrupt); err != nil {
			return err
		}
	}

	if len(p.options.DownBinPath) == 0 {
		return nil
	}

	p.downCmd = exec.Command(p.options.DownBinPath, p.options.DownArgs...)

	if err := runner.Outputs(p.downCmd); err != nil {
		return err
	}

	if err := p.downCmd.Run(); err != nil {
		return err
	}

	return nil
}

func (p *customProcess) String() string {
	return "custom"
}

func NewProcess(opts ...runner.ProcessOption) runner.Process {
	options := runner.NewProcessOptions(opts...)

	p := &customProcess{
		options: options,
		mtx:     sync.RWMutex{},
	}

	return p
}
