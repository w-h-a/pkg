package sidecar

import (
	"time"
)

type Event struct {
	EventName string      `json:"eventName,omitempty"`
	To        []string    `json:"to,omitempty"`
	CreatedAt time.Time   `json:"createdAt,omitempty"`
	State     State       `json:"state,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type State struct {
	StoreId string   `json:"storeId,omitempty"`
	Records []Record `json:"records,omitempty"`
}

type Record struct {
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}
