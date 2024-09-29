package log

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var (
	level  Level      = LevelInfo
	prefix string     = ""
	format FormatFunc = func(r Record) string {
		b, _ := json.Marshal(r)
		return string(b)
	}
)

var logger Log

type Log interface {
	Options() LogOptions
	Write(Record) error
	Read(opts ...ReadOption) ([]Record, error)
	String() string
}

type FormatFunc func(Record) string

func SetLogger(l Log) {
	logger = l
}

func GetLogger() Log {
	return logger
}

// Trace provides trace level logging
func Trace(v ...interface{}) {
	WithLevel(LevelTrace, v...)
}

// Tracef provides trace level logging
func Tracef(format string, v ...interface{}) {
	WithLevelf(LevelTrace, format, v...)
}

// Debug provides debug level logging
func Debug(v ...interface{}) {
	WithLevel(LevelDebug, v...)
}

// Debugf provides debug level logging
func Debugf(format string, v ...interface{}) {
	WithLevelf(LevelDebug, format, v...)
}

// Warn provides warn level logging
func Warn(v ...interface{}) {
	WithLevel(LevelWarn, v...)
}

// Warnf provides warn level logging
func Warnf(format string, v ...interface{}) {
	WithLevelf(LevelWarn, format, v...)
}

// Info provides info level logging
func Info(v ...interface{}) {
	WithLevel(LevelInfo, v...)
}

// Infof provides info level logging
func Infof(format string, v ...interface{}) {
	WithLevelf(LevelInfo, format, v...)
}

// Error provides warn level logging
func Error(v ...interface{}) {
	WithLevel(LevelError, v...)
}

// Errorf provides warn level logging
func Errorf(format string, v ...interface{}) {
	WithLevelf(LevelError, format, v...)
}

// Fatal logs with Log and then exits with os.Exit(1)
func Fatal(v ...interface{}) {
	WithLevel(LevelFatal, v...)
	os.Exit(1)
}

// Fatalf logs with Logf and then exits with os.Exit(1)
func Fatalf(format string, v ...interface{}) {
	WithLevelf(LevelFatal, format, v...)
	os.Exit(1)
}

// WithLevel logs with the level specified
func WithLevel(l Level, v ...interface{}) {
	if l > logger.Options().Level {
		return
	}
	log(l, v...)
}

// WithLevel logs with the level specified
func WithLevelf(l Level, format string, v ...interface{}) {
	if l > logger.Options().Level {
		return
	}
	logf(l, format, v...)
}

func log(l Level, v ...interface{}) {
	if len(logger.Options().Prefix) > 0 {
		v = append([]interface{}{logger.Options().Prefix, " "}, v...)
	}

	logger.Write(
		Record{
			Timestamp: time.Now(),
			Message:   fmt.Sprint(v...),
			Metadata:  map[string]string{"level": l.String()},
		},
	)
}

func logf(l Level, format string, v ...interface{}) {
	if len(logger.Options().Prefix) > 0 {
		format = logger.Options().Prefix + " " + format
	}

	logger.Write(
		Record{
			Timestamp: time.Now(),
			Message:   fmt.Sprintf(format, v...),
			Metadata:  map[string]string{"level": l.String()},
		},
	)
}
