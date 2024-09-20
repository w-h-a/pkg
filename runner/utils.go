package runner

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ioPair struct {
	in  io.ReadCloser
	out *os.File
}

func Outputs(cmd *exec.Cmd) error {
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
						fmt.Fprintf(out, "%s", s)
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

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port, nil
}

type ParallelTest struct {
	funs []func(c *assert.CollectT)
	mtx  sync.RWMutex
}

func (p *ParallelTest) Add(fun func(c *assert.CollectT)) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.funs = append(p.funs, fun)
}

func (p *ParallelTest) worker(jobCh <-chan func()) {
	for job := range jobCh {
		job()
	}
}

func NewParallelTest(t *testing.T, funs ...func(c *assert.CollectT)) *ParallelTest {
	p := &ParallelTest{
		funs: funs,
	}

	t.Cleanup(func() {
		p.mtx.Lock()
		defer p.mtx.Unlock()

		t.Helper()

		jobCh := make(chan func(), len(p.funs))

		workerCount := 4

		for i := 0; i < workerCount; i++ {
			go p.worker(jobCh)
		}

		cs := make([]*assert.CollectT, len(p.funs))

		var wg sync.WaitGroup

		wg.Add(len(p.funs))

		for i := range p.funs {
			cs[i] = &assert.CollectT{}

			jobCh <- func() {
				defer wg.Done()
				defer recover()

				p.funs[i](cs[i])
			}
		}

		wg.Wait()

		close(jobCh)
	})

	return p
}
