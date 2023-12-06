package postgrestor

import "github.com/Schalure/gofermart/internal/configs"

type Storage struct {
}

func NewStorage(config *configs.Config) *Storage {

	return &Storage{}
}