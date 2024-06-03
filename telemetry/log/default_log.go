package log

import (
	golog "log"

	"github.com/w-h-a/pkg/telemetry/buffer"
	"github.com/w-h-a/pkg/telemetry/buffer/memory"
)

type defaultLog struct {
	options LogOptions
	buffer  buffer.Buffer
}

func (l *defaultLog) Options() LogOptions {
	return l.options
}

func (l *defaultLog) Write(r Record) error {
	out := l.options.Format(r)
	golog.Print(out)
	l.buffer.Put(out)
	return nil
}

func (l *defaultLog) Read(opts ...ReadOption) ([]Record, error) {
	options := NewReadOptions(opts...)

	entries := []*buffer.Entry{}

	if options.Count > 0 {
		entries = l.buffer.Get(options.Count)
	}

	records := []Record{}

	for _, entry := range entries {
		record := Record{
			Message:   entry.Value,
			Timestamp: entry.Timestamp,
		}

		records = append(records, record)
	}

	return records, nil
}

func (l *defaultLog) String() string {
	return "default"
}

func NewLog(opts ...LogOption) Log {
	options := NewLogOptions(opts...)

	l := &defaultLog{
		options: options,
		buffer:  memory.NewBuffer(buffer.BufferWithSize(options.Size)),
	}

	return l
}
