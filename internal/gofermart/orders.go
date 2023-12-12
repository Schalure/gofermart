package gofermart

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/Schalure/gofermart/internal/gofermart/gofermaterrors"
	"github.com/Schalure/gofermart/internal/storage"
)


func (g *Gofermart) LoadOrder(login, orderNumber string) error {

	pc := "func (g *Gofermart) LoadOrder(login, orderNumber string) error"

	if !g.isOrderValid(orderNumber) {
		g.loggerer.Infow(
			pc,
			"message", "order number is not valid",
			"orderNumber", orderNumber,
		)
		return gofermaterrors.InvalidOrderNumber
	}

	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second * 5)
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
			return gofermaterrors.DublicateOrderNumberByUser
		} else if order.UserLogin != login {
			return gofermaterrors.DublicateOrderNumber
		}
		return gofermaterrors.Internal
	}

	order = storage.Order{
		OrderNumber: orderNumber,
		UserLogin: login,
		OrderStatus: storage.OrderStatusNew,
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel2()
	err = g.storager.AddNewOrder(ctx2, order)
	if err != nil {
		g.loggerer.Infow(
			pc,
			"message", "can't add new order",
			"orderNumber", orderNumber,
			"user", login,
			"error", err,
		)
		return gofermaterrors.Internal
	}

	err = g.orderProvider.addOrderToCash(order)
	if err != nil {
		g.loggerer.Debugw(
			pc,
			"message", err,
		)
	}
	return nil
}

func (g *Gofermart) checkOrderStatusWorker(ctx context.Context) {

	jobsCh := make(chan storage.Order)
	resultCh := make(chan storage.Order)

	for {
		if freeCels, ok := g.orderProvider.isCashNotFull(); ok {
			ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
			orders, err :=  g.storager.GetOrdersToUpdateStatus(ctx, freeCels)
			cancel()
			if err != nil {
				break
			}
			for _, order := range orders {
				if err := g.orderProvider.addOrderToCash(order); err != nil {
					break
				}
			}
		}
	}
}

func (g *Gofermart) GetOrders(login string) ([]storage.Order, error) {


	return nil, nil
}

//	Order number validity check
func (g *Gofermart) isOrderValid(orderNumber string) bool {

	if !g.validOrderNumber.MatchString(orderNumber) {
		return false
	}
	
	orderNumberStringArr := strings.Split(orderNumber, "")
	if len(orderNumberStringArr) < 2 {
		return false
	}
	orderNumberArr := make([]int, len(orderNumberStringArr))
	for i, s := range orderNumberStringArr {

		orderNumberArr[i], _ = strconv.Atoi(s)
	}

	return LunaAlgorithm(orderNumberArr)
}

//	Checking by the Luna algorithm
func LunaAlgorithm(data []int) bool {

	const k = 9
	dataLenght := len(data)
	startNum := dataLenght % 2
	res := 0

	for ; startNum < dataLenght - 1; startNum += 2 {
		data[startNum] *= 2
		if data[startNum] >= k {
			data[startNum] -= k
		}
		res += data[startNum] + data[startNum + 1]
	}

	return res % 10 == 0
}
