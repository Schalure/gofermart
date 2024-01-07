package gofermart

import (
	"context"
	"errors"
	"time"

	"github.com/Schalure/gofermart/internal/storage"
)

func (g *Gofermart) Withdraw(ctx context.Context, login, orderNumber string, sum float64) error {

	pc := "func (g *Gofermart) Withdraw(login, order string, sum int) error"

	if !g.isOrderValid(orderNumber) {
		g.loggerer.Infow(
			pc,
			"message", "order number is not valid",
			"orderNumber", orderNumber,
		)
		return ErrInvalidOrderNumber
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*1000)
	err := g.storager.WithdrawPointsForOrder(ctx, login, orderNumber, sum, time.Now())
	cancel()
	if err != nil {
		g.loggerer.Infow(
			"can't withdraw points for order",
			"loggin", login,
			"orderNumber", orderNumber,
			"sum", sum,
			"error", err,
		)

		if errors.Is(err, storage.ErrInsufficientFunds) {
			return ErrInsufficientFunds
		}
		if errors.Is(err, storage.ErrOrderNumberAlreadyExists) {
			return ErrDublicateOrderNumber
		}
		return ErrInternal
	}
	return nil
}

func (g *Gofermart) GetWithdraws(ctx context.Context, login string) ([]storage.Order, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	orders, err := g.storager.GetPointWithdraws(ctx, login)
	cancel()
	if err != nil {
		return nil, ErrNoOrdersForPoints
	}
	return orders, nil
}
