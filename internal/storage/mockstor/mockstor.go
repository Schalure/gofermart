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

func (s *Storage) AddNewOrder(ctx context.Context, order storage.Order) error {

	panic("no implemented: func (s *Storage) AddNewOrder(ctx context.Context, order storage.Order) error")
}

func (s *Storage) GetOrderByNumber(ctx context.Context, orderNumber string) (storage.Order, error) {

	panic("no implemented: func (s *Storage) GetOrderByNumber(ctx context.Context, orderNumber string) (storage.Order, error)")
}

func (s *Storage) GetOrdersByLogin(ctx context.Context, login string) ([]storage.Order, error) {

	panic("no implemented: func (s *Storage) GetOrdersByLogin(ctx context.Context, login string) ([]storage.Order, error)")
}

func (s *Storage) GetOrdersToUpdateStatus(ctx context.Context) ([]storage.Order, error) {

	panic("no implemented: func (s *Storage) GetOrdersToUpdateStatus(ctx context.Context, maxCount int) ([]storage.Order, error)")
}


func (s *Storage) WithdrawPointsForOrder(ctx context.Context, orderNumber string, sum int) error {

	panic("no implemented: func (s *Storage) WithdrawPointsForOrder(ctx context.Context, orderNumber string, sum int) error")

}
func (s *Storage) GetPointWithdraws(ctx context.Context, login string) ([]storage.Order, error) {

	panic("no implemented: func (s *Storage) GetPointWithdraws(ctx context.Context, login string) ([]storage.Order, error)")

}