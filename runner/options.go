package runner

import (
	"context"

	"github.com/google/uuid"
)

type RunnerOption func(o *RunnerOptions)

type RunnerOptions struct {
	Id      string
	Files   []File
	Context context.Context
}

func RunnerWithId(id string) RunnerOption {
	return func(o *RunnerOptions) {
		o.Id = id
	}
}

func RunnerWithFiles(files ...File) RunnerOption {
	return func(o *RunnerOptions) {
		o.Files = files
	}
}

func NewRunnerOptions(opts ...RunnerOption) RunnerOptions {
	options := RunnerOptions{
		Id:      uuid.New().String(),
		Files:   []File{},
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type ManagerOption func(o *ManagerOptions)

type ManagerOptions struct {
	File    File
	Context context.Context
}

func ManagerWithFile(file File) ManagerOption {
	return func(o *ManagerOptions) {
		o.File = file
	}
}

func NewManagerOptions(opts ...ManagerOption) ManagerOptions {
	options := ManagerOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
