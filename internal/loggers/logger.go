package loggers

import (
	"fmt"
	"log/slog"

	"github.com/Schalure/gofermart/internal/gofermart"
)

type LoggerType string

func (l LoggerType) String() string {
	return fmt.Sprintf("Logger type: %s", l)
}

const (
	Slog LoggerType = "Slog"
	Zap LoggerType = "Zap"
)


func NewLogger(loggerType LoggerType) gofermart.Logger {

	switch loggerType {
	case Slog:
	}
	return &slog.Logger{}
}