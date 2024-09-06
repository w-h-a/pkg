package docker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/w-h-a/pkg/runner"
)

type dockerManager struct {
	options runner.ManagerOptions
}

type ioPair struct {
	in  io.ReadCloser
	out *os.File
}

func (d *dockerManager) Options() runner.ManagerOptions {
	return d.options
}

func (d *dockerManager) Apply() error {
	cmd := exec.Command("docker", "compose", "--file", d.options.File.Path, "up", "--build", "--detach")

	if err := d.outputs(cmd); err != nil {
		return err
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (d *dockerManager) Destroy() error {
	cmd := exec.Command("docker", "compose", "--file", d.options.File.Path, "down", "--volumes")

	if err := d.outputs(cmd); err != nil {
		return err
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (d *dockerManager) String() string {
	return "docker"
}

func (*dockerManager) outputs(cmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	pairs := []ioPair{
		{in: stdout, out: os.Stdout},
		{in: stderr, out: os.Stderr},
	}

	for _, ioPair := range pairs {
		go func(in io.ReadCloser, out *os.File) {
			defer in.Close()

			reader := bufio.NewReader(in)

			for {
				s, err := reader.ReadString('\n')
				if err == nil || err == io.EOF {
					if len(strings.TrimSpace(s)) != 0 {
						fmt.Fprintf(out, "%s\n", s)
					}
					if err == io.EOF {
						return
					}
				} else {
					fmt.Fprintf(out, "error: %s\n", err.Error())
					return
				}
			}
		}(ioPair.in, ioPair.out)
	}
	
	return nil
}

func NewManager(opts ...runner.ManagerOption) runner.Manager {
	options := runner.NewManagerOptions(opts...)

	d := &dockerManager{
		options: options,
	}

	return d
}
