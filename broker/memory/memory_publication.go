package memory

import "github.com/w-h-a/pkg/broker"

type memoryPublication struct {
	topic   string
	message *broker.Message
}

func (p *memoryPublication) Topic() string {
	return p.topic
}

func (p *memoryPublication) Message() *broker.Message {
	return p.message
}

func (p *memoryPublication) Ack() error {
	return nil
}

func (p *memoryPublication) String() string {
	return "memory"
}
