package grpcstreams

import (
	"fmt"
	"time"

	"github.com/w-h-a/pkg/streams"
)

func SendEvent(sub streams.Subscriber, event *streams.Event) error {
	cpy := *event

	if sub.Options().AutoAck {
		sub.Channel() <- cpy
		return nil
	}

	cpy.SetAck(Ack(sub, cpy))
	cpy.SetNack(Nack(sub, cpy))

	sub.SetAttempts(0, cpy)

	tick := time.NewTicker(sub.Options().AckWait)
	defer tick.Stop()

	for range tick.C {
		count, ok := sub.GetAttempts(cpy)
		if !ok {
			break
		}

		if sub.Options().RetryLimit > -1 && count > sub.Options().RetryLimit {
			sub.Ack(cpy)
			return fmt.Errorf("discarding event %s because the number of attempts %d exceeded the retry limit %d", cpy.Id, count, sub.Options().RetryLimit)
		}

		sub.Channel() <- cpy

		sub.SetAttempts(count+1, cpy)
	}

	return nil
}

func Ack(sub streams.Subscriber, event streams.Event) func() error {
	return func() error {
		return sub.Ack(event)
	}
}

func Nack(sub streams.Subscriber, event streams.Event) func() error {
	return func() error {
		return sub.Nack(event)
	}
}
