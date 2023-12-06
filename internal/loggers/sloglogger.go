package loggers

import (
	"log/slog"
	"os"

	"github.com/Schalure/gofermart/internal/configs"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(config *configs.Config) *Logger {

	var loggerLevel slog.Level

	switch config.AppConfig.Env {
	case configs.Debug:
		loggerLevel = slog.LevelDebug
	case configs.Local:
		loggerLevel = slog.LevelInfo
	case configs.Prod:
		loggerLevel = slog.LevelError
	}

	loggerHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggerLevel,
	})

	logger := slog.New(loggerHandler)

	return &Logger{
		logger: logger,
	}
}