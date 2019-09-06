package log

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-godin/log/level"
	"github.com/go-kit/kit/log"
	stdzipkin "github.com/openzipkin/zipkin-go"
)

// Logger is the default logging interface which is used throughout all godin services.
type Logger interface {
	Log(keyvals ...interface{})
	Debug(message string, keyvals ...interface{})
	Info(message string, keyvals ...interface{})
	Warning(message string, keyvals ...interface{})
	Error(message string, keyvals ...interface{})
	With(keyvals ...interface{}) Log
	WithTrace(ctx context.Context) Log
}

const (
	LevelDebug          = "debug"
	LevelInfo           = "info"
	LevelWarning        = "warning"
	LevelError          = "error"
	MessageKey          = "message"
	EnvironmentVariable = "LOG_LEVEL"
)

type Log struct {
	kitLogger log.Logger
	span      stdzipkin.Span
}

// NewLogger creates a new, leveled Log. The given level is the allowed minimal level.
func NewLogger(logLevel string) Log {
	levelOpt, err := evaluateLogLevel(logLevel)

	var kitLogger log.Logger
	kitLogger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	kitLogger = level.NewFilter(kitLogger, levelOpt)

	log := Log{
		kitLogger: kitLogger,
	}

	// the error from evaluateLogLevel needs to be logged
	if err != nil {
		log.Warning("", "err", err)
	}

	return log
}

// NewLoggerFromEnv creates a new Log, configuring the log level using an environment variable.
func NewLoggerFromEnv() Log {
	levelStr := os.Getenv(EnvironmentVariable)
	return NewLogger(levelStr)
}

func (l Log) SetLevel(logLevel string) {
	lvl, err := evaluateLogLevel(logLevel)
	if err != nil {
		lvl = level.AllowInfo()
	}
	l.kitLogger = level.NewFilter(l.kitLogger, lvl)
}

func (l Log) WithTrace(ctx context.Context) Log {
	if span := stdzipkin.SpanFromContext(ctx); span != nil {
		return Log{
			kitLogger: l.kitLogger,
			span:      span,
		}
	}
	return Log{
		kitLogger: l.kitLogger,
		span:      nil,
	}
}

// Log redirects to go-kit/log.Log
func (l Log) Log(keyvals ...interface{}) {
	l.handleTrace("", keyvals)
	_ = l.kitLogger.Log(keyvals...)
}

// Debug will log a message and arbitrary key-value pairs
func (l Log) Debug(message string, keyvals ...interface{}) {
	_ = level.Debug(l.kitLogger).Log(l.mergeKeyValues(message, keyvals)...)
}

// Info will log a message and arbitrary key-value pairs
func (l Log) Info(message string, keyvals ...interface{}) {
	l.handleTrace(message, keyvals)
	_ = level.Info(l.kitLogger).Log(l.mergeKeyValues(message, keyvals)...)
}

// Warning will log a message and arbitrary key-value pairs
func (l Log) Warning(message string, keyvals ...interface{}) {
	l.handleTrace(message, keyvals)
	_ = level.Warn(l.kitLogger).Log(l.mergeKeyValues(message, keyvals)...)
}

// Error will log a message and arbitrary key-value pairs
func (l Log) Error(message string, keyvals ...interface{}) {
	l.handleTrace(message, keyvals)
	_ = level.Error(l.kitLogger).Log(l.mergeKeyValues(message, keyvals)...)
}

func (l Log) With(keyvals ...interface{}) Log {
	if len(keyvals) == 0 {
		return l
	}

	kitLogger := log.With(l.kitLogger, keyvals...)

	return Log{
		kitLogger: kitLogger,
		span:      l.span,
	}
}

func (l Log) handleTrace(message string, keyvals []interface{}) {
	if l.span != nil {
		if message != "" {
			l.span.Annotate(time.Now(), message)
		}
		for i := 0; i < len(keyvals); i += 2 {
			if i >= len(keyvals) || i+1 >= len(keyvals) {
				break // break only for the uneven keyval combination, all others will be tagged
			}
			l.span.Tag(fmt.Sprint(keyvals[i]), fmt.Sprint(keyvals[i+1]))
		}
	}
}

// evaluateLogLevel maps a given logLevel as string (e.g. from an ENV variable) to a level Option.
// If the passed logLevel does not exist, all levels will be enabled by default.
func evaluateLogLevel(logLevel string) (level.Option, error) {
	logLevel = strings.ToLower(logLevel)
	switch logLevel {
	case LevelDebug:
		return level.AllowDebug(), nil
	case LevelInfo:
		return level.AllowInfo(), nil
	case LevelWarning:
		return level.AllowWarn(), nil
	case LevelError:
		return level.AllowError(), nil
	default:
		return level.AllowAll(), fmt.Errorf("no log-level passed, falling back to debug")
	}
}

// mergeKeyValues will append the level and message field to already existing keyvals
func (l Log) mergeKeyValues(message string, keyvals []interface{}) []interface{} {
	var list []interface{}
	var levelData []interface{}

	if message != "" {
		levelData = append(levelData, MessageKey)
		levelData = append(levelData, message)
	}

	list = append(list, levelData...)
	list = append(list, keyvals...)

	return list
}
