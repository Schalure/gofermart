package gofermart

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Schalure/gofermart/internal/gofermart/gofermaterrors"
	"github.com/Schalure/gofermart/internal/storage"
	"github.com/jackc/pgx/pgtype"
)

// Add new order to system
func (g *Gofermart) LoadOrder(ctx context.Context, login, orderNumber string) error {

	pc := "func (g *Gofermart) LoadOrder(login, orderNumber string) error"

	if !g.isOrderValid(orderNumber) {
		g.loggerer.Infow(
			pc,
			"message", "order number is not valid",
			"orderNumber", orderNumber,
		)
		return gofermaterrors.InvalidOrderNumber
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
			return gofermaterrors.DublicateOrderNumberByUser
		}
		if order.UserLogin != login {
			return gofermaterrors.DublicateOrderNumber
		}
		return gofermaterrors.Internal
	}

	order = storage.Order{
		OrderNumber: orderNumber,
		UserLogin:   login,
		OrderStatus: storage.OrderStatusNew,
		UploadedOrder: pgtype.Timestamptz{Time: time.Now()},
	}

	ctx2, cancel2 := context.WithTimeout(ctx, time.Second*5)
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

	g.wg.Add(1)
	go g.addToInputCh(order)

	return nil
}

// Return info about orders by user
func (g *Gofermart) GetOrders(ctx context.Context, login string) ([]storage.Order, error) {

	pc := "func (g *Gofermart) GetOrders(login string) ([]storage.Order, error)"

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	orders, err := g.storager.GetOrdersByLogin(ctx, login)
	if err != nil {
		g.loggerer.Infow(
			pc,
			"message", "error reading user orders",
			"user", login,
			"error", err,
		)
		return nil, err
	}
	return orders, nil
}

func (g *Gofermart) orderCheckWorker(ctx context.Context) {

	//	1.	get orders to check from database
	ctxGetOrders, cancelGetOrders := context.WithTimeout(ctx, time.Second * 5)
	orders, err := g.storager.GetOrdersToUpdateStatus(ctxGetOrders)
	cancelGetOrders()
	if err != nil {
		return
	}

	for _, order := range orders {
		g.wg.Add(1)
		go g.addToInputCh(order)
	}

	//	2.	run workers
	resultCh := make(chan int)
	for i := 0; i < numWorkers; i++ {
		go g.worker(ctx, i, resultCh)
	}

	//	3.	run orderCheckWorker task
	for {
		select {
		case <-ctx.Done():
			g.doneCh <- struct{}{}
			close(g.doneCh)
			g.wg.Wait()
			close(g.inputCh)
			close(resultCh)
			return
		case status := <- resultCh:
			if status == http.StatusTooManyRequests {
				t := time.NewTimer(time.Second * 60)
				select {
				case <-ctx.Done():
					g.doneCh <- struct{}{}
					close(g.doneCh)
					g.wg.Wait()
					close(g.inputCh)
					close(resultCh)
					return
				case <-t.C:
				}
			}
		}
	}
}

func (g *Gofermart) worker(ctx context.Context, workerID int, resultCh chan<- int) {
	
	pc := "func (g *Gofermart) worker(ctx context.Context, chanID int, resultCh chan<- int)"

	for order := range g.inputCh {
		g.loggerer.Debugw(
			pc,
			"message", "start processing order",
			"workerId", workerID,
			"order", order.String(),
		)
		ctx1, cancel1 := context.WithTimeout(ctx, time.Second)
		res, status := g.orderChecker.OrderCheck(ctx1, order.OrderNumber)
		cancel1()

		if status == http.StatusOK && res.OrderStatus != order.OrderStatus {
			ctx2, cancel2 := context.WithTimeout(ctx, time.Second)
			cancel2()
			err := g.storager.UpdateOrder(ctx2, res.OrderNumber, res.OrderStatus, res.BonusPoints)
			if err != nil {
				g.wg.Add(1)
				go g.addToInputCh(order)
				continue
			}
		}

		if res.OrderStatus == storage.OrderStatusNew || res.OrderStatus == storage.OrderStatusProcessing {
			g.wg.Add(1)
			go g.addToInputCh(order)
		}

		resultCh <- status
		g.loggerer.Debugw(
			pc,
			"message", "finish processing order",
			"workerId", workerID,
			"order", order.String(),
			"status", status,
		)
	}
}

//	Add order to input channel
func (g *Gofermart) addToInputCh(order storage.Order) {

	defer g.wg.Done()
	select {
	case <- g.doneCh:
		return
	case g.inputCh <- order:
	}
}

// Order number validity check
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

// Checking by the Luna algorithm
func LunaAlgorithm(data []int) bool {

	const k = 9
	dataLenght := len(data)
	startNum := dataLenght % 2
	res := 0

	for ; startNum < dataLenght-1; startNum += 2 {
		data[startNum] *= 2
		if data[startNum] >= k {
			data[startNum] -= k
		}
		res += data[startNum] + data[startNum+1]
	}

	return res%10 == 0
}
