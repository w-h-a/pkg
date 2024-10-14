package tracev2

import (
	"time"
)

type SpanData struct {
	Name     string            `json:"name"`
	Id       string            `json:"id"`
	Parent   string            `json:"parent"`
	Trace    string            `json:"trace"`
	Started  time.Time         `json:"started"`
	Ended    time.Time         `json:"ended"`
	Metadata map[string]string `json:"metadata"`
	Status   Status            `json:"status"`
}

type Status struct {
	Code        uint32 `json:"code"`
	Description string `json:"description"`
}
