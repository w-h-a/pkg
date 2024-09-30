package buffer

import "context"

type BufferOption func(o *BufferOptions)

type BufferOptions struct {
	Size    int
	Context context.Context
}

func BufferWithSize(s int) BufferOption {
	return func(o *BufferOptions) {
		o.Size = s
	}
}

func NewBufferOptions(opts ...BufferOption) BufferOptions {
	options := BufferOptions{
		Context: context.Background(),
		Size:    defaultSize,
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
