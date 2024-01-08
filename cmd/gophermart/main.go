package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Schalure/gofermart/internal/configs"
	"github.com/Schalure/gofermart/internal/gofermart"
	"github.com/Schalure/gofermart/internal/loggers"
	"github.com/Schalure/gofermart/internal/loyaltysystem"
	"github.com/Schalure/gofermart/internal/server"
	"github.com/Schalure/gofermart/internal/storage/postgrestor"
)

func main() {

	log.Println("Starting application initialization...")

	ctx, cancel := context.WithCancel(context.Background())

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
	var wgStop sync.WaitGroup
	wgStop.Add(1)
	defer wgStop.Wait()
	service.Run(ctx, &wgStop)

	log.Println("HTTP server initializing...")
	handler := server.NewHandler(service, service, service)
	midleware := server.NewMidleware(logger, service)
	server := server.NewServer(config.EnvConfig.ServiceHost, handler, midleware)

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		log.Println("Application stoped by os.Interript...")

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Stop(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		shutdownCancel()
		cancel()
		log.Println("server shutdowned...")
	}()

	log.Println(config.EnvConfig.ServiceHost, config.EnvConfig.DBHost, config.EnvConfig.AccrualHost)
	log.Println("Gofermart service have been started...")

	err = server.Run()
	log.Println("Gofermart service stop:", err)
	<-ctx.Done()
}
