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
		return gofermaterrors.InvalidOrderNumber
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	order, err := g.storager.GetOrderByNumber(ctx, orderNumber)
	cancel()
	if err != nil {
		return gofermaterrors.InvalidOrderNumber
	}

	if order.UserLogin != login {
		return gofermaterrors.InvalidOrderNumber
	}

	ctx, cancel = context.WithTimeout(ctx, time.Second*5)
	user, err := g.storager.GetUserByLogin(ctx, login)
	cancel()
	if err != nil {
		return err
	}

	if user.LoyaltyPoints < sum {
		return gofermaterrors.InsufficientFunds
	}

	user.LoyaltyPoints -= sum
	user.WithdrawnPoints += sum

	ctx, cancel = context.WithTimeout(ctx, time.Second*5)
	err = g.storager.WithdrawPointsForOrder(ctx, orderNumber, sum, time.Now())
	cancel()
	if err != nil {
		return gofermaterrors.Internal
	}

	return nil
}

func (g *Gofermart) GetWithdraws(ctx context.Context, login string) ([]storage.Order, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	orders, err := g.storager.GetPointWithdraws(ctx, login)
	cancel()
	if err != nil {
		return nil, gofermaterrors.NoOrdersForPoints
	}
	return orders, nil
}
