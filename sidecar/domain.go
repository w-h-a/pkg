package sidecar

type Event struct {
	EventName string      `json:"eventName,omitempty"`
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

type Secret struct {
	Data map[string]string `json:"data,omitempty"`
}
