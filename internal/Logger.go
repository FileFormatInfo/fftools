package internal

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

var Logger = initLogger()
var LogLevel = slog.LevelInfo

// level=Trace is for stuff that should not be logged in production, either because it's too verbose or because it contains sensitive information.

func initLogger() *LoggerWithTrace {

	logLevel := os.Getenv("LOG_LEVEL")
	var err error
	if logLevel != "" {
		ilvl, atoiErr := strconv.Atoi(logLevel)
		if atoiErr == nil {
			LogLevel = slog.Level(ilvl)
		} else if strings.ToLower(logLevel) == "trace" {
			LogLevel = -8
		} else {
			err = LogLevel.UnmarshalText([]byte(logLevel))
			if err != nil {
				LogLevel = slog.LevelInfo
			}
		}
	}

	// get log format from the environment
	logFormat := os.Getenv("LOG_FORMAT")
	var handler slog.Handler
	switch logFormat {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level:     LogLevel,
			AddSource: true,
		})
	default:
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level:     LogLevel,
			AddSource: true,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	if err != nil {
		slog.Error("unable to set log level", "level", LogLevel, "error", err)
	} else {
		slog.Debug("log level set", "level", LogLevel)
	}

	return &LoggerWithTrace{Logger: logger}
}

// LoggerWithTrace is a wrapper around slog.Logger that adds a Trace method.
type LoggerWithTrace struct {
	*slog.Logger
}

// Trace logs a message at the trace level.
func (l *LoggerWithTrace) Trace(msg string, keyvals ...interface{}) {
	l.Log(context.Background(), slog.Level(-8), msg, keyvals...)
}

// WithTrace returns the LoggerWithTrace instance.
func (l *LoggerWithTrace) WithTrace() *LoggerWithTrace {
	return l
}
