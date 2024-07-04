package streams

import (
	"encoding/json"
	"errors"
	"time"
)

type Ack func() error

type Nack func() error

type Event struct {
	Id        string            `json:"id"`
	Topic     string            `json:"topic"`
	Payload   []byte            `json:"payload"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata"`
	Processed map[string]bool   `json:"processed"`
	ack       Ack
	nack      Nack
}

func (e *Event) Unmarshal(v interface{}) error {
	return json.Unmarshal(e.Payload, v)
}

func (e *Event) SetAck(f Ack) {
	e.ack = f
}

func (e *Event) Ack() error {
	if e.ack == nil {
		return errors.New("no ack function set")
	}
	return e.ack()
}

func (e *Event) SetNack(f Nack) {
	e.nack = f
}

func (e *Event) Nack() error {
	if e.nack == nil {
		return errors.New("no nack function set")
	}
	return e.nack()
}
