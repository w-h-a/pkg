package log

import "time"

const (
	LevelFatal Level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

type Level int32

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelWarn:
		return "warn"
	case LevelInfo:
		return "info"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

type Record struct {
	Message   interface{}       `json:"message"`
	Timestamp time.Time         `json:"time"`
	Metadata  map[string]string `json:"metadata"`
}
