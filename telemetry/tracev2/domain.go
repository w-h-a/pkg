package tracev2

import (
	"time"
)

// TODO: status
type SpanData struct {
	Name     string
	Id       string
	Parent   string
	Trace    string
	Started  time.Time
	Ended    time.Time
	Metadata map[string]string
}
