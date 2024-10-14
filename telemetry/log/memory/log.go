package memory

import (
	golog "log"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/utils/memoryutils"
)

type memoryLog struct {
	options log.LogOptions
	buffer  *memoryutils.Buffer
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

	var entries []*memoryutils.Entry

	if options.Count > 0 {
		entries = l.buffer.Get(options.Count)
	} else {
		entries = l.buffer.Get(l.buffer.Options().Size)
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
		l.buffer = memoryutils.NewBuffer(memoryutils.BufferWithSize(s))
	} else {
		l.buffer = memoryutils.NewBuffer()
	}

	return l
}
