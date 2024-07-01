package streams

import (
	"encoding/json"
	"errors"
	"time"
)

type AckFunc func() error

type NackFunc func() error

type Event struct {
	Id        string
	Topic     string
	Payload   []byte
	Timestamp time.Time
	Metadata  map[string]string
	ackFunc   AckFunc
	nackFunc  NackFunc
}

func (e *Event) Unmarshal(v interface{}) error {
	return json.Unmarshal(e.Payload, v)
}

func (e *Event) SetAck(f AckFunc) {
	e.ackFunc = f
}

func (e *Event) Ack() error {
	if e.ackFunc == nil {
		return errors.New("no ack function set")
	}
	return e.ackFunc()
}

func (e *Event) SetNack(f NackFunc) {
	e.nackFunc = f
}

func (e *Event) Nack() error {
	if e.nackFunc == nil {
		return errors.New("no nack function set")
	}
	return e.nackFunc()
}
