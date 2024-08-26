package sidecar

import (
	"github.com/w-h-a/pkg/store"
)

type Sidecar interface {
	Options() SidecarOptions
	OnEventPublished(event *Event) error
	SaveStateToStore(state *State) error
	RetrieveStateFromStore(store, key string) ([]*store.Record, error)
	ReadEventsFromBroker(broker string)
	UnsubscribeFromBroker(broker string) error
	String() string
}
