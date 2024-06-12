package store

import "time"

type Record struct {
	Key    string
	Value  []byte
	Expiry time.Duration
}
