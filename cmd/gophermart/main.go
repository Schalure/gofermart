package main

import (
	"log"

	"github.com/Schalure/gofermart/internal/configs"
	"github.com/Schalure/gofermart/internal/loggers"
)

func main() {

	log.Println("Starting application initialization...")

	log.Println("Config initializing...")
	config := configs.NewConfig()

	log.Println("Logger initializing...")
	logger := loggers.NewLogger(loggers.Slog)
}
