package nats

import "github.com/w-h-a/pkg/broker"

type natsPublication struct {
	topic   string
	message *broker.Message
}

func (p *natsPublication) Topic() string {
	return p.topic
}

func (p *natsPublication) Message() *broker.Message {
	return p.message
}

func (p *natsPublication) Ack() error {
	return nil
}

func (p *natsPublication) String() string {
	return "nats"
}
