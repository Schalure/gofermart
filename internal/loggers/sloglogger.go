package loggers

import (
	"log/slog"
	"os"

	"github.com/Schalure/gofermart/internal/configs"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(environment configs.Environment) *Logger {

	var loggerLevel slog.Level

	switch environment {
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

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info("", args...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, keysAndValues...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debug(msg, keysAndValues...)
}
