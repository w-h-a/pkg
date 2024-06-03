package log

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var (
	logger Log        = NewLog()
	level  Level      = LevelInfo
	prefix string     = ""
	size   int        = 1024
	format FormatFunc = func(r Record) string {
		b, _ := json.Marshal(r)
		return string(b)
	}
)

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

func SetLevel(l Level) {
	level = l
}

func GetLevel() Level {
	return level
}

func SetPrefix(p string) {
	prefix = p
}

func SetServiceName(name string) {
	prefix = fmt.Sprintf("[%s]", name)
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
	if l > level {
		return
	}
	log(v...)
}

// WithLevel logs with the level specified
func WithLevelf(l Level, format string, v ...interface{}) {
	if l > level {
		return
	}
	logf(format, v...)
}

func log(v ...interface{}) {
	if len(prefix) > 0 {
		v = append([]interface{}{prefix, " "}, v...)
	}

	logger.Write(
		Record{
			Timestamp: time.Now(),
			Message:   fmt.Sprint(v...),
			Metadata:  map[string]string{"level": level.String()},
		},
	)
}

func logf(format string, v ...interface{}) {
	if len(prefix) > 0 {
		format = prefix + " " + format
	}

	logger.Write(
		Record{
			Timestamp: time.Now(),
			Message:   fmt.Sprintf(format, v...),
			Metadata:  map[string]string{"level": level.String()},
		},
	)
}
