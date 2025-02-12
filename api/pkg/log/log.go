package hablog

import (
	"log/slog"
	"os"
	"strings"
)

var (
	logOptions = struct {
		level   slog.Level
		handler string
	}{
		level:   slog.LevelInfo,
		handler: "default",
	}
)

func GetLogger() *slog.Logger {
	var handler slog.Handler
	switch logOptions.handler {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     logOptions.level,
		})
	default:
		handler = newHandler(os.Stderr)
	}

	return slog.New(handler)
}

func Init(logLevel string, logHandler string) {
	if logLevel != "" {
		logLevel = strings.ToLower(logLevel)
	}
	switch logLevel {
	case "debug":
		logOptions.level = slog.LevelDebug
	case "info":
		logOptions.level = slog.LevelInfo
	case "warn":
		logOptions.level = slog.LevelWarn
	case "error":
		logOptions.level = slog.LevelError
	default:
		logOptions.level = slog.LevelInfo
	}

	if logHandler == "" {
		logHandler = "default"
	}
	logOptions.handler = strings.ToLower(logHandler)
}

func init() {
	logLevel := os.Getenv("AFB_LOG_LEVEL")
	logHandler := os.Getenv("AFB_LOG_HANDLER")
	Init(logLevel, logHandler)
}
