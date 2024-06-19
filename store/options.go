package store

import (
	"context"
)

type StoreOption func(o *StoreOptions)

type StoreOptions struct {
	Nodes    []string
	Database string
	Table    string
	Context  context.Context
}

func StoreWithNodes(addrs ...string) StoreOption {
	return func(o *StoreOptions) {
		o.Nodes = addrs
	}
}

func StoreWithDatabase(db string) StoreOption {
	return func(o *StoreOptions) {
		o.Database = db
	}
}

func StoreWithTable(tbl string) StoreOption {
	return func(o *StoreOptions) {
		o.Table = tbl
	}
}

func NewStoreOptions(opts ...StoreOption) StoreOptions {
	options := StoreOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type WriteOption func(o *WriteOptions)

type WriteOptions struct{}

func NewWriteOptions(opts ...WriteOption) WriteOptions {
	options := WriteOptions{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type ReadOption func(o *ReadOptions)

type ReadOptions struct {
	Prefix bool
	Suffix bool
	Limit  uint
	Offset uint
}

func ReadWithPrefix() ReadOption {
	return func(o *ReadOptions) {
		o.Prefix = true
	}
}

func ReadWithSuffix() ReadOption {
	return func(o *ReadOptions) {
		o.Suffix = true
	}
}

func ReadWithLimit(lim uint) ReadOption {
	return func(o *ReadOptions) {
		o.Limit = lim
	}
}

func ReadWithOffset(off uint) ReadOption {
	return func(o *ReadOptions) {
		o.Offset = off
	}
}

func NewReadOptions(opts ...ReadOption) ReadOptions {
	options := ReadOptions{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type ListOption func(o *ListOptions)

type ListOptions struct {
	Prefix string
	Suffix string
	Limit  uint
	Offset uint
}

func ListWithPrefix(p string) ListOption {
	return func(o *ListOptions) {
		o.Prefix = p
	}
}

func ListWithSuffix(s string) ListOption {
	return func(o *ListOptions) {
		o.Suffix = s
	}
}

func ListWithLimit(lim uint) ListOption {
	return func(o *ListOptions) {
		o.Limit = lim
	}
}

func ListWithOffset(off uint) ListOption {
	return func(o *ListOptions) {
		o.Offset = off
	}
}

func NewListOptions(opts ...ListOption) ListOptions {
	options := ListOptions{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}

type DeleteOption func(o *DeleteOptions)

type DeleteOptions struct{}

func NewDeleteOptions(opts ...DeleteOption) DeleteOptions {
	options := DeleteOptions{}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
