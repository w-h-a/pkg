package memoryutils

import (
	"sync"
	"time"
)

var (
	defaultSize = 1024
)

type Buffer struct {
	options BufferOptions
	mtx     sync.RWMutex
	entries []*Entry
}

type Entry struct {
	Value     interface{}
	Timestamp time.Time
}

func (m *Buffer) Options() BufferOptions {
	return m.options
}

func (m *Buffer) Put(v interface{}) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	// make the entry
	entry := &Entry{
		Value:     v,
		Timestamp: time.Now(),
	}

	// append the entry
	m.entries = append(m.entries, entry)

	// if the length of the entries is greater than the
	// specified size, then trim down the buffer by 1
	if len(m.entries) > m.options.Size {
		m.entries = m.entries[1:]
	}
}

func (m *Buffer) Get(n int) []*Entry {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	// reset bad inputs
	if n > len(m.entries) || n < 0 {
		n = len(m.entries)
	}

	// create a delta
	delta := len(m.entries) - n

	return m.entries[delta:]
}

func NewBuffer(opts ...BufferOption) *Buffer {
	options := NewBufferOptions(opts...)

	b := &Buffer{
		options: options,
		mtx:     sync.RWMutex{},
		entries: []*Entry{},
	}

	return b
}
