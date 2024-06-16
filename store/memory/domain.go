package memory

import "time"

type InternalRecord struct {
	Key       string
	Value     []byte
	ExpiresAt time.Time
}
