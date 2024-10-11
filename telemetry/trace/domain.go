package trace

import "time"

type Span struct {
	Name     string            `json:"name"`
	Id       string            `json:"id"`
	Parent   string            `json:"parent"`
	Trace    string            `json:"trace"`
	Started  time.Time         `json:"started"`
	Duration time.Duration     `json:"duration"`
	Metadata map[string]string `json:"metadata"`
}
