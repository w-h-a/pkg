package store

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Store interface {
	Options() StoreOptions
	Write(rec *Record, opts ...WriteOption) error
	Read(key string, opts ...ReadOption) ([]*Record, error)
	List(opts ...ListOption) ([]string, error)
	Delete(key string, opts ...DeleteOption) error
	String() string
}
