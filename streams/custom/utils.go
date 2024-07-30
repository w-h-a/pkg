package custom

import (
	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/streams"
)

func Ack(s client.Stream, sub streams.Subscriber, event streams.Event) func() error {
	return func() error {
		return s.Send(sub.Ack(event))
	}
}

func Nack(s client.Stream, sub streams.Subscriber, event streams.Event) func() error {
	return func() error {
		return s.Send(sub.Nack(event))
	}
}
