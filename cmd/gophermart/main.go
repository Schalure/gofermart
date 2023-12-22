package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/Schalure/gofermart/internal/configs"
	"github.com/Schalure/gofermart/internal/gofermart"
	"github.com/Schalure/gofermart/internal/loggers"
	"github.com/Schalure/gofermart/internal/loyaltysystem"
	"github.com/Schalure/gofermart/internal/server"
	"github.com/Schalure/gofermart/internal/storage/postgrestor"
)

func main() {

	log.Println("Starting application initialization...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("Config initializing...")
	config, err := configs.NewConfig()
	if err != nil {
		log.Println("Config have been initialized with error:", err)
	}

	log.Println("Logger initializing...")
	logger := loggers.NewLogger(config.AppConfig.Env)

	log.Println("Storage initializing...")
	storage, err := postgrestor.NewStorage(config.EnvConfig.DBHost)
	if err != nil {
		log.Panicln("Storage have been initialized with error:", err)
	}

	log.Println("Service initializing...")
	orderChecker := loyaltysystem.NewLoyaltySystem(config.EnvConfig.AccrualHost)
	service := gofermart.NewGofermart(
		storage,
		logger,
		orderChecker,
		config.AppConfig.LoginRules,
		config.AppConfig.PassRules,
		config.AppConfig.OrderNumberRules,
		config.AppConfig.TokenTTL,
	)
	service.Run(ctx)

	log.Println("HTTP server initializing...")
	handler := server.NewHandler(service, service, service)
	midleware := server.NewMidleware(logger, service)
	server := server.NewServer(handler, midleware)

	log.Println("Gofermart service have been started...")
	err = server.Run(config.EnvConfig.ServiceHost)
	server.Stop(err)
}
