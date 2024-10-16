package sidecar

import (
	"context"
	"errors"

	"github.com/w-h-a/pkg/store"
)

var (
	ErrComponentNotFound = errors.New("component not found")
	ErrInvalidGroupName  = errors.New("subscriber group name should be of form <group>-<topic>")
)

type Sidecar interface {
	Options() SidecarOptions
	SaveStateToStore(ctx context.Context, state *State) error
	ListStateFromStore(ctx context.Context, store string) ([]*store.Record, error)
	SingleStateFromStore(ctx context.Context, store, key string) ([]*store.Record, error)
	RemoveStateFromStore(ctx context.Context, store, key string) error
	WriteEventToBroker(ctx context.Context, event *Event) error
	ReadEventsFromBroker(ctx context.Context, broker string)
	UnsubscribeFromBroker(ctx context.Context, broker string) error
	ReadFromSecretStore(ctx context.Context, secretStore string, name string) (*Secret, error)
	String() string
}
