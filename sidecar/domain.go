package sidecar

import (
	"time"

	"github.com/w-h-a/pkg/store"
)

type Event struct {
	EventName string      `json:"eventName,omitempty"`
	To        []string    `json:"to,omitempty"`
	CreatedAt time.Time   `json:"createdAt,omitempty"`
	State     *State      `json:"state,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type State struct {
	StoreId string          `json:"storeId,omitempty"`
	Records []*store.Record `json:"records,omitempty"`
}
