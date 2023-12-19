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

	//	Create table users
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXIST users(
		login VARCHAR(64) PRIMARY KEY,
		password VARCHAR(64) NOT NULL,
		loyalty_points MONEY DEFAULT(0),
		withdrawn_points MONEY DEFAULT(0));
	`)
	if err != nil {
		return nil, err
	}

	//	Create order status constants
	_, err = db.Exec(ctx, 
		`CREATE TYPE order_state AS ENUM ($1, $2, $3, $4);`, 
		storage.OrderStatusNew, storage.OrderStatusProcessing, storage.OrderStatusInvalid, storage.OrderStatusProcessed) 
		if err != nil {
			return nil, err
		}

	_, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXIST orders(
		order_number VARCHAR(64) PRIMARY KEY,
		order_status order_state DEFAULT($1),
		uploaded_order TIMESTAMP WITH TIME ZONE NOT NULL,
		bonus_points MONEY DEFAULT(0),
		uploaded_bonus TIMESTAMP WITH TIME ZONE,
		login VARCHAR(64)
		FOREIGN KEY (login) REFERENCES users(login) ON DELETE CASCADE
	);`, storage.OrderStatusNew)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) AddNewUser(ctx context.Context, user storage.User) error {

	_, err := s.db.Exec(ctx, `INSERT INTO users (login, password) VALUES($1, $2);`, user.Login, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUserByLogin(ctx context.Context, login string) (storage.User, error) {

	var  user storage.User
	row := s.db.QueryRow(ctx, `SELECT * FROM users WHERE login = $1;`, login)
	err := row.Scan(&user.Login, &user.Password, &user.LoyaltyPoints, &user.WithdrawnPoints)
	return user, err
}

func (s *Storage) AddNewOrder(ctx context.Context, order storage.Order) error {

	_, err := s.db.Exec(ctx,
		`INSERT INTO orders (order_number, order_status, uploaded_order, login) VALUES($1, $2, $3, $4);`,
		order.OrderNumber, order.OrderStatus, order.UploadedOrder.Format(time.RFC3339), order.UserLogin,
	)
	if err != nil {
		return err
	}
	return nil	
}

func (s *Storage) GetOrderByNumber(ctx context.Context, orderNumber string) (storage.Order, error) {

	var order storage.Order
	row := s.db.QueryRow(ctx, `SELECT * FROM orders WHERE order_number = $1;`, orderNumber)
	err := row.Scan(&order.OrderNumber, &order.OrderStatus, &order.UploadedOrder, &order.BonusPoints, &order.UploadedBonus, &order.UserLogin)
	return order, err
}

func (s *Storage) GetOrdersByLogin(ctx context.Context, login string) ([]storage.Order, error) {

	orders := make([]storage.Order, 0)

	rows, err := s.db.Query(ctx,
		`SELECT * FROM orders WHERE login = $1 ORDER BY uploaded_order DESC;`,
		login,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o storage.Order
		err := rows.Scan(&o.OrderNumber, &o.OrderStatus, &o.UploadedOrder, &o.BonusPoints, &o.UploadedBonus, &o.UserLogin)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (s *Storage) GetOrdersToUpdateStatus(ctx context.Context) ([]storage.Order, error) {

	orders := make([]storage.Order, 0)

	rows, err := s.db.Query(ctx,
		`SELECT * FROM orders WHERE (order_status = $1 OR order_status = $2) ORDER BY uploaded_order;`,
		storage.OrderStatusNew, storage.OrderStatusProcessing,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o storage.Order
		err := rows.Scan(&o.OrderNumber, &o.OrderStatus, &o.UploadedOrder, &o.BonusPoints, &o.UploadedBonus, &o.UserLogin)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (s *Storage) WithdrawPointsForOrder(ctx context.Context, orderNumber string, sum float64, uploaded_at time.Time) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx,
		`UPDATE users SET bonus_points = bonus_points - $1, uploaded_bonus = uploaded_bonus + $2
		WHERE login = (SELECT login FROM orders WHERE order_number = $3);`,
		sum, sum, orderNumber,
	); err != nil {
		return err
	}

	if _, err = tx.Exec(ctx,
		`UPDATE orders SET (bonus_points = $1, uploaded_bonus = $2) WHERE order_number = $3;`,
		sum, uploaded_at, orderNumber,
	); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Storage) GetPointWithdraws(ctx context.Context, login string) ([]storage.Order, error) {

	orders := make([]storage.Order, 0)

	rows, err := s.db.Query(ctx,
		`SELECT * FROM orders WHERE login = $1 ORDER BY uploaded_bonus DESC;`,
		login,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o storage.Order
		err := rows.Scan(&o.OrderNumber, &o.OrderStatus, &o.UploadedOrder, &o.BonusPoints, &o.UploadedBonus, &o.UserLogin)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
