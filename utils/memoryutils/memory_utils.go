package memoryutils

import (
	"sync"
	"time"
)

type Buffer struct {
	Size    int
	mtx     sync.RWMutex
	entries []*Entry
}

type Entry struct {
	Value     interface{}
	Timestamp time.Time
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
	if len(m.entries) > m.Size {
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

func NewBuffer(size int) *Buffer {
	b := &Buffer{
		Size:    size,
		mtx:     sync.RWMutex{},
		entries: []*Entry{},
	}

	return b
}
