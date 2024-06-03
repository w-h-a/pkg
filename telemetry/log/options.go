package log

import "context"

type LogOption func(o *LogOptions)

type LogOptions struct {
	Format  FormatFunc
	Size    int
	Context context.Context
}

func LogWithFormat(f FormatFunc) LogOption {
	return func(o *LogOptions) {
		o.Format = f
	}
}

func LogWithSize(s int) LogOption {
	return func(o *LogOptions) {
		o.Size = s
	}
}

func NewLogOptions(opts ...LogOption) LogOptions {
	options := LogOptions{
		Context: context.Background(),
		Format:  format,
		Size:    size,
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type ReadOption func(o *ReadOptions)

type ReadOptions struct {
	Count int
}

func ReadWithCount(c int) ReadOption {
	return func(o *ReadOptions) {
		o.Count = c
	}
}

func NewReadOptions(opts ...ReadOption) ReadOptions {
	options := ReadOptions{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
