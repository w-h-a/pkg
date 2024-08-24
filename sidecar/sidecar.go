package sidecar

import (
	"github.com/google/uuid"
	"github.com/w-h-a/pkg/store"
)

var (
	defaultID = uuid.New().String()
)

type Sidecar interface {
	Options() SidecarOptions
	OnEventPublished(event *Event) error
	SaveStateToStore(store string, state []*store.Record) error
	RetrieveStateFromStore(store, key string) ([]*store.Record, error)
	ReadEventsFromBroker(broker, eventName string)
	UnsubscribeFromBroker(broker string) error
	String() string
}
