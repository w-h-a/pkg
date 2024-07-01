package streams

type ProduceOption func(o *ProduceOptions)

type ProduceOptions struct {
}

func NewProduceOptions(opts ...ProduceOption) ProduceOptions {
	options := ProduceOptions{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type ConsumeOption func(o *ConsumeOptions)

type ConsumeOptions struct {
}

func NewConsumeOptions(opts ...ConsumeOption) ConsumeOptions {
	options := ConsumeOptions{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
