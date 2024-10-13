package tracev2

import (
	"time"
)

// TODO: status
type SpanData struct {
	Name     string            `json:"name"`
	Id       string            `json:"id"`
	Parent   string            `json:"parent"`
	Trace    string            `json:"trace"`
	Started  time.Time         `json:"started"`
	Ended    time.Time         `json:"ended"`
	Metadata map[string]string `json:"metadata"`
}
