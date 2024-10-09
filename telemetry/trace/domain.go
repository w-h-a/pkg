package trace

import "time"

type Span struct {
	Name     string
	Id       string
	Parent   string
	Trace    string
	Started  time.Time
	Duration time.Duration
	Metadata map[string]string
}
