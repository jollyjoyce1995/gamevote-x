package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	once sync.Once
)

// Init initializes the global slog logger.
// It uses JSON handler if LOG_FORMAT=json, otherwise Text handler.
// Default level is Info, can be overridden with LOG_LEVEL.
func Init() {
	once.Do(func() {
		var handler slog.Handler
		opts := &slog.HandlerOptions{
			Level: getLogLevel(),
		}

		format := os.Getenv("LOG_FORMAT")
		if format == "json" {
			handler = slog.NewJSONHandler(os.Stdout, opts)
		} else {
			handler = slog.NewTextHandler(os.Stdout, opts)
		}

		logger := slog.New(handler)
		slog.SetDefault(logger)
	})
}

func getLogLevel() slog.Level {
	levelStr := os.Getenv("LOG_LEVEL")
	switch levelStr {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
