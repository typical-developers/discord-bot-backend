package logger

import (
	"log/slog"
	"os"

	"github.com/dusted-go/logging/prettylog"
)

type Logger struct {
	*slog.Logger
	WithSource *slog.Logger
}

var Log Logger

func init() {
	logLevel := slog.LevelInfo

	if os.Getenv("ENVIRONMENT") == "development" {
		logLevel = slog.LevelDebug
	}

	Log = Logger{
		Logger: slog.New(prettylog.NewHandler(&slog.HandlerOptions{
			Level:       logLevel,
			AddSource:   false,
			ReplaceAttr: nil,
		})),
		WithSource: slog.New(prettylog.NewHandler(&slog.HandlerOptions{
			Level:       logLevel,
			AddSource:   true,
			ReplaceAttr: nil,
		})),
	}
}
