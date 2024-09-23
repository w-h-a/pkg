package sidecar

import (
	"errors"

	"github.com/w-h-a/pkg/store"
)

var (
	ErrComponentNotFound = errors.New("component not found")
	ErrInvalidGroupName  = errors.New("subscriber group name should be of form <group>-<topic>")
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
	ReadFromSecretStore(secretStore string, name, prefix string) (*Secret, error)
	String() string
}
