package memory

import (
	golog "log"

	"github.com/w-h-a/pkg/telemetry/buffer"
	"github.com/w-h-a/pkg/telemetry/buffer/memory"
	"github.com/w-h-a/pkg/telemetry/log"
)

// TODO: put this in its own memory pkg
type defaultLog struct {
	options log.LogOptions
	buffer  buffer.Buffer
}

func (l *defaultLog) Options() log.LogOptions {
	return l.options
}

func (l *defaultLog) Write(r log.Record) error {
	out := l.options.Format(r)
	golog.Print(out)
	l.buffer.Put(out)
	return nil
}

func (l *defaultLog) Read(opts ...log.ReadOption) ([]log.Record, error) {
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

func (l *defaultLog) String() string {
	return "default"
}

func NewLog(opts ...log.LogOption) log.Log {
	options := log.NewLogOptions(opts...)

	l := &defaultLog{
		options: options,
	}

	if s, ok := GetSizeFromContext(options.Context); ok && s > 0 {
		l.buffer = memory.NewBuffer(buffer.BufferWithSize(s))
	} else {
		l.buffer = memory.NewBuffer()
	}

	log.SetLogger(l)

	return l
}
