package runner

import (
	"context"

	"github.com/google/uuid"
)

type RunnerOption func(o *RunnerOptions)

type RunnerOptions struct {
	Id        string
	Processes []Process
	Context   context.Context
}

func RunnerWithId(id string) RunnerOption {
	return func(o *RunnerOptions) {
		o.Id = id
	}
}

func RunnerWithProcesses(ps ...Process) RunnerOption {
	return func(o *RunnerOptions) {
		o.Processes = ps
	}
}

func NewRunnerOptions(opts ...RunnerOption) RunnerOptions {
	options := RunnerOptions{
		Id:        uuid.New().String(),
		Processes: []Process{},
		Context:   context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type ProcessOption func(o *ProcessOptions)

type ProcessOptions struct {
	Id          string
	UpBinPath   string
	UpArgs      []string
	DownBinPath string
	DownArgs    []string
	EnvVars     map[string]string
	Context     context.Context
}

func ProcessWithId(id string) ProcessOption {
	return func(o *ProcessOptions) {
		o.Id = id
	}
}

func ProcessWithUpBinPath(path string) ProcessOption {
	return func(o *ProcessOptions) {
		o.UpBinPath = path
	}
}

func ProcessWithUpArgs(args ...string) ProcessOption {
	return func(o *ProcessOptions) {
		o.UpArgs = args
	}
}

func ProcessWithDownBinPath(path string) ProcessOption {
	return func(o *ProcessOptions) {
		o.DownBinPath = path
	}
}

func ProcessWithDownArgs(args ...string) ProcessOption {
	return func(o *ProcessOptions) {
		o.DownArgs = args
	}
}

func ProcessWithEnvVars(envs map[string]string) ProcessOption {
	return func(o *ProcessOptions) {
		o.EnvVars = envs
	}
}

func NewProcessOptions(opts ...ProcessOption) ProcessOptions {
	options := ProcessOptions{
		UpArgs:   []string{},
		DownArgs: []string{},
		EnvVars:  map[string]string{},
		Context:  context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
