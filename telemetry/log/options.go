package log

import (
	"context"
	"fmt"
)

type LogOption func(o *LogOptions)

type LogOptions struct {
	Prefix  string
	Format  FormatFunc
	Context context.Context
}

func LogWithPrefix(prefix string) LogOption {
	return func(o *LogOptions) {
		o.Prefix = fmt.Sprintf("[%s]", prefix)
	}
}

func LogWithFormat(f FormatFunc) LogOption {
	return func(o *LogOptions) {
		o.Format = f
	}
}

func NewLogOptions(opts ...LogOption) LogOptions {
	options := LogOptions{
		Prefix:  prefix,
		Format:  format,
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type ReadOption func(o *ReadOptions)

type ReadOptions struct {
	Count   int
	Context context.Context
}

func ReadWithCount(c int) ReadOption {
	return func(o *ReadOptions) {
		o.Count = c
	}
}

func NewReadOptions(opts ...ReadOption) ReadOptions {
	options := ReadOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
