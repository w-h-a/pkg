package sidecar

import (
	"github.com/w-h-a/pkg/store"
)

type Sidecar interface {
	Options() SidecarOptions
	SaveStateToStore(state *State) error
	RetrieveStateFromStore(store, key string) ([]*store.Record, error)
	RemoveStateFromStore(store, key string) error
	WriteEventToBroker(event *Event) error
	ReadEventsFromBroker(broker string)
	UnsubscribeFromBroker(broker string) error
	String() string
}
