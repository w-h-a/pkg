package sidecar

import (
	"time"
)

type Event struct {
	EventName  string      `json:"eventName,omitempty"`
	To         []string    `json:"to,omitempty"`
	Concurrent string      `json:"concurrent,omitempty"`
	State      State       `json:"state,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	CreatedAt  time.Time   `json:"createdAt,omitempty"`
}

type State struct {
	StoreId string   `json:"storeId,omitempty"`
	Records []Record `json:"records,omitempty"`
}

type Record struct {
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}
