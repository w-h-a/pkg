package streams

import (
	"encoding/json"
	"errors"
	"time"
)

type Ack func() error

type Nack func() error

type Event struct {
	Id        string
	Topic     string
	Payload   []byte
	Timestamp time.Time
	Metadata  map[string]string
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
