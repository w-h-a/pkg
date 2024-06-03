package memory

import (
	"sync"
	"time"

	"github.com/w-h-a/pkg/telemetry/buffer"
)

type memoryBuffer struct {
	options buffer.BufferOptions
	mtx     sync.RWMutex
	entries []*buffer.Entry
}

func (m *memoryBuffer) Options() buffer.BufferOptions {
	return m.options
}

func (m *memoryBuffer) Put(v interface{}) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	// make the entry
	entry := &buffer.Entry{
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

func (m *memoryBuffer) Get(n int) []*buffer.Entry {
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

func (m *memoryBuffer) String() string {
	return "memory"
}

func NewBuffer(opts ...buffer.BufferOption) buffer.Buffer {
	options := buffer.NewBufferOptions(opts...)

	b := &memoryBuffer{
		options: options,
		mtx:     sync.RWMutex{},
		entries: []*buffer.Entry{},
	}

	return b
}
