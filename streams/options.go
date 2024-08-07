package streams

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/w-h-a/pkg/store"
)

type StreamsOption func(o *StreamsOptions)

type StreamsOptions struct {
	Store   store.Store
	Context context.Context
}

func StreamsWithStore(s store.Store) StreamsOption {
	return func(o *StreamsOptions) {
		o.Store = s
	}
}

func NewStreamsOptions(opts ...StreamsOption) StreamsOptions {
	options := StreamsOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type SubscribeOption func(o *SubscribeOptions)

type SubscribeOptions struct {
	Group      string
	Topic      string
	AckWait    time.Duration
	RetryLimit int
	Offset     time.Time
}

func SubscribeWithGroup(n string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Group = n
	}
}

func SubscribeWithTopic(t string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Topic = t
	}
}

func SubscribeWithAckWait(ackWait time.Duration) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AckWait = ackWait
	}
}

func SubscribeWithRetryLimit(retries int) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.RetryLimit = retries
	}
}

func SubscribeWithOffset(t time.Time) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Offset = t
	}
}

func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	options := SubscribeOptions{
		Group:      uuid.New().String(),
		AckWait:    16 * time.Second,
		RetryLimit: 4,
		Offset:     time.Now().Add(time.Hour * -1),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type ProduceOption func(o *ProduceOptions)

type ProduceOptions struct {
	Timestamp time.Time
	Metadata  map[string]string
	Processed map[string]bool
	Context   context.Context
}

func ProduceWithTimestamp(t time.Time) ProduceOption {
	return func(o *ProduceOptions) {
		o.Timestamp = t
	}
}

func ProduceWithMetadata(md map[string]string) ProduceOption {
	return func(o *ProduceOptions) {
		o.Metadata = md
	}
}

func ProduceWithProcessed(p map[string]bool) ProduceOption {
	return func(o *ProduceOptions) {
		o.Processed = p
	}
}

func ProduceWithContext(ctx context.Context) ProduceOption {
	return func(o *ProduceOptions) {
		o.Context = ctx
	}
}

func NewProduceOptions(opts ...ProduceOption) ProduceOptions {
	options := ProduceOptions{
		Timestamp: time.Now(),
		Metadata:  map[string]string{},
		Processed: map[string]bool{},
		Context:   context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
