package logger

import (
	"log/slog"
	"os"

	"github.com/dusted-go/logging/prettylog"
)

var Log *slog.Logger

func init() {
	logLevel := slog.LevelInfo

	if os.Getenv("ENVIRONMENT") == "development" {
		logLevel = slog.LevelDebug
	}

	prettyHandler := prettylog.NewHandler(&slog.HandlerOptions{
		Level:       logLevel,
		AddSource:   true,
		ReplaceAttr: nil,
	})

	Log = slog.New(prettyHandler)
}
