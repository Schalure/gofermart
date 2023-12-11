package postgrestor

import (
	"context"

	"github.com/Schalure/gofermart/internal/storage"
)

type Storage struct {
}

func NewStorage() *Storage {

	return &Storage{}
}

func (s *Storage) AddNewUser(ctx context.Context, user storage.User) error {

	panic("no implemented: func (s *Storage) AddNewUser(ctx context.Context, user storage.User) error")
}

func (s *Storage) GetUserByLogin(ctx context.Context, login string) (storage.User, error) {

	panic("no implemented: func (s *Storage) GetUserByLogin(ctx context.Context, login string) (storage.User, error)")

}