package memoryutils

type BufferOption func(o *BufferOptions)

type BufferOptions struct {
	Size int
}

func BufferWithSize(size int) BufferOption {
	return func(o *BufferOptions) {
		o.Size = size
	}
}

func NewBufferOptions(opts ...BufferOption) BufferOptions {
	options := BufferOptions{
		Size: defaultSize,
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
