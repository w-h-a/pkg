package sidecar

import (
	"errors"

	"github.com/w-h-a/pkg/store"
)

var (
	ErrComponentNotFound = errors.New("component not found")
)

type Sidecar interface {
	Options() SidecarOptions
	SaveStateToStore(state *State) error
	ListStateFromStore(store string) ([]*store.Record, error)
	SingleStateFromStore(store, key string) ([]*store.Record, error)
	RemoveStateFromStore(store, key string) error
	WriteEventToBroker(event *Event) error
	ReadEventsFromBroker(broker string)
	UnsubscribeFromBroker(broker string) error
	String() string
}
