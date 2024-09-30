package memory

import (
	golog "log"

	"github.com/w-h-a/pkg/telemetry/buffer"
	"github.com/w-h-a/pkg/telemetry/buffer/memory"
	"github.com/w-h-a/pkg/telemetry/log"
)

type memoryLog struct {
	options log.LogOptions
	buffer  buffer.Buffer
}

func (l *memoryLog) Options() log.LogOptions {
	return l.options
}

func (l *memoryLog) Write(rec log.Record) error {
	out := l.options.Format(rec)

	golog.Print(out)

	l.buffer.Put(out)

	return nil
}

func (l *memoryLog) Read(opts ...log.ReadOption) ([]log.Record, error) {
	options := log.NewReadOptions(opts...)

	entries := []*buffer.Entry{}

	if options.Count > 0 {
		entries = l.buffer.Get(options.Count)
	}

	records := []log.Record{}

	for _, entry := range entries {
		record := log.Record{
			Message:   entry.Value,
			Timestamp: entry.Timestamp,
		}

		records = append(records, record)
	}

	return records, nil
}

func (l *memoryLog) String() string {
	return "memory"
}

func NewLog(opts ...log.LogOption) log.Log {
	options := log.NewLogOptions(opts...)

	l := &memoryLog{
		options: options,
	}

	if s, ok := GetSizeFromContext(options.Context); ok && s > 0 {
		l.buffer = memory.NewBuffer(buffer.BufferWithSize(s))
	} else {
		l.buffer = memory.NewBuffer()
	}

	return l
}
