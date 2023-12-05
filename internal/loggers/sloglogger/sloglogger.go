package sloglogger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func (l *Logger) NewLogger() *Logger {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &Logger{
		logger: logger,
	}
}