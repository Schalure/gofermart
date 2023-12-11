package mockstor

import (
	"context"
	"fmt"

	"github.com/Schalure/gofermart/internal/storage"
)

type Storage struct {
	Users map[string]storage.User
}

func NewStorage() *Storage {

	return &Storage{
		Users: make(map[string]storage.User),
	}
}

func (s *Storage) AddNewUser(ctx context.Context, user storage.User) error {

	s.Users[user.Login] = user
	return nil
}

func (s *Storage) GetUserByLogin(ctx context.Context, login string) (storage.User, error) {

	var err error
	user, ok := s.Users[login]
	if !ok {
		err = fmt.Errorf("user not found")
	}
	return user, err
}
