package postgrestor

import (
	"context"
	"time"

	"github.com/Schalure/gofermart/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(dbConnectionString string) (*Storage, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db, err := pgxpool.New(ctx, dbConnectionString)
	if err != nil {
		return nil, err
	}

	return &Storage{}, nil
}

func (s *Storage) AddNewUser(ctx context.Context, user storage.User) error {

	panic("no implemented: func (s *Storage) AddNewUser(ctx context.Context, user storage.User) error")
}

func (s *Storage) GetUserByLogin(ctx context.Context, login string) (storage.User, error) {

	panic("no implemented: func (s *Storage) GetUserByLogin(ctx context.Context, login string) (storage.User, error)")

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

func (s *Storage) WithdrawPointsForOrder(ctx context.Context, orderNumber string, sum float64) error {

	panic("no implemented: func (s *Storage) WithdrawPointsForOrder(ctx context.Context, orderNumber string, sum int) error")

}
func (s *Storage) GetPointWithdraws(ctx context.Context, login string) ([]storage.Order, error) {

	panic("no implemented: func (s *Storage) GetPointWithdraws(ctx context.Context, login string) ([]storage.Order, error)")

}
