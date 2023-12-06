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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &Logger{
		logger: logger,
	}
}