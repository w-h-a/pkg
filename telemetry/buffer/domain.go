package buffer

import "time"

type Entry struct {
	Value     interface{}
	Timestamp time.Time
}
