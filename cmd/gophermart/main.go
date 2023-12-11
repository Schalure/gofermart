package main

import (
	"log"

	"github.com/Schalure/gofermart/internal/configs"
	"github.com/Schalure/gofermart/internal/gofermart"
	"github.com/Schalure/gofermart/internal/loggers"
	"github.com/Schalure/gofermart/internal/server"
	"github.com/Schalure/gofermart/internal/storage/postgrestor"
)

func main() {

	log.Println("Starting application initialization...")

	log.Println("Config initializing...")
	config, err := configs.NewConfig()
	if err != nil {
		log.Println("Config have been initialized with error:", err)
	}

	log.Println("Logger initializing...")
	logger := loggers.NewLogger(config.AppConfig.Env)

	log.Println("Storage initializing...")
	storage := postgrestor.NewStorage()
	
	log.Println("Service initializing...")
	service := gofermart.NewGofermart(storage, logger, config.AppConfig.LoginRules, config.AppConfig.PassRules, config.AppConfig.TokenTTL)

	log.Println("HTTP server initializing...")
	handler := server.NewHandler(service)
	midleware := server.NewMidleware(logger, service)
	server := server.NewServer(handler, midleware)

	log.Println("Gofermart service have been started...")
	err = server.Run(config.EnvConfig.ServiceHost)
	server.Stop(err)
}
