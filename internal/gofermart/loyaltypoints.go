package gofermart

import (
	"context"
	"time"

	"github.com/Schalure/gofermart/internal/gofermart/gofermaterrors"
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
		return gofermaterrors.ErrInvalidOrderNumber
	}

	ctx1, cancel1 := context.WithTimeout(ctx, time.Second*5)
	defer cancel1()
	order, err := g.storager.GetOrderByNumber(ctx1, orderNumber)
	if err == nil {
		g.loggerer.Infow(
			pc,
			"message", "can't get order by number",
			"orderNumber", orderNumber,
			"error", err,
		)
		if order.UserLogin == login {
			return gofermaterrors.ErrDublicateOrderNumberByUser
		}
		if order.UserLogin != login {
			return gofermaterrors.ErrDublicateOrderNumber
		}
		return gofermaterrors.ErrInternal
	}

	ctx2, cancel2 := context.WithTimeout(ctx, time.Second*5)
	user, err := g.storager.GetUserByLogin(ctx2, login)
	cancel2()
	if err != nil {
		return gofermaterrors.ErrInternal
	}

	if user.LoyaltyPoints < sum {
		return gofermaterrors.ErrInsufficientFunds
	}

	user.LoyaltyPoints -= sum
	user.WithdrawnPoints += sum

	ctx3, cancel3 := context.WithTimeout(ctx, time.Second*5)
	err = g.storager.WithdrawPointsForOrder(ctx3, login, orderNumber, sum, time.Now())
	cancel3()
	if err != nil {
		return gofermaterrors.ErrInternal
	}

	return nil
}

func (g *Gofermart) GetWithdraws(ctx context.Context, login string) ([]storage.Order, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	orders, err := g.storager.GetPointWithdraws(ctx, login)
	cancel()
	if err != nil {
		return nil, gofermaterrors.ErrNoOrdersForPoints
	}
	return orders, nil
}
