package server

import (
	"testing"

	"github.com/Schalure/gofermart/internal/configs"
	"github.com/Schalure/gofermart/internal/gofermart"
	"github.com/Schalure/gofermart/internal/loggers"
	"github.com/Schalure/gofermart/internal/server/midlewares"
	"github.com/Schalure/gofermart/internal/storage/mockstor"
)

func Test_UserRegistration(t *testing.T) {

	
	config, _ := configs.NewConfig()
	logger := loggers.NewLogger(config)
	stor := mockstor.NewStorage()
	service := gofermart.NewGofermart(config, stor, logger)

	server := NewServer(config, NewHandler(service), midlewares.NewMidleware())	
	
}