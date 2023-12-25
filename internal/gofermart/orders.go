package gofermart

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Schalure/gofermart/internal/gofermart/gofermaterrors"
	"github.com/Schalure/gofermart/internal/storage"
	"github.com/jackc/pgx/pgtype"
)

const RepetitiveCheckTime = 0//time.Second * 5
const SleepCheckTime = 0//time.Second * 60

var count int

// Add new order to system
func (g *Gofermart) LoadOrder(ctx context.Context, login, orderNumber string) error {

	pc := "func (g *Gofermart) LoadOrder(login, orderNumber string) error"

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

	order = storage.Order{
		OrderNumber:   orderNumber,
		UserLogin:     login,
		OrderStatus:   storage.OrderStatusNew,
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
		return gofermaterrors.ErrInternal
	}

	g.wg.Add(1)
	go g.addToInputCh(order, 0)

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

	var wgWait sync.WaitGroup
	defer wgWait.Wait()

	//	1.	get orders to check from database
	ctxGetOrders, cancelGetOrders := context.WithTimeout(ctx, time.Second*5)
	orders, err := g.storager.GetOrdersToUpdateStatus(ctxGetOrders)
	cancelGetOrders()
	if err != nil {
		g.loggerer.Debugw("orderCheckWorker err with read database",
			"err", err,
		)
		return
	}

	for _, order := range orders {
		g.wg.Add(1)
		go g.addToInputCh(order, 0)
	}

	//	2.	run workers
	resultCh := make(chan int)
	pauseSignalCh := make(chan struct{}, numWorkers)
	for i := 0; i < numWorkers; i++ {
		wgWait.Add(1)
		go g.worker(ctx, &wgWait, i, resultCh, pauseSignalCh)
	}

	g.loggerer.Debugw("orderCheckWorker start",)
	//	3.	run orderCheckWorker task
	for {
		select {
		case <-ctx.Done():
			g.loggerer.Debugw("количество запущеных горутин", "count", count)
			g.doneCh <- struct{}{}
			close(g.doneCh)
			g.loggerer.Debugw("closed g.doneCh")
			g.wg.Wait()
			close(g.inputCh)
			g.loggerer.Debugw("closed g.inputCh")
			close(resultCh)
			g.loggerer.Debugw("closed resultCh")
			close(pauseSignalCh)
			g.loggerer.Debugw("closed pauseSignalCh")
			return
		case status := <-resultCh:
			g.loggerer.Debugw("return status from resultCh",
				"status", status,
			)
			if status == http.StatusTooManyRequests {
				g.loggerer.Debugw("paused orderCheckWorker")
				for i := 0; i < numWorkers; i++ {
					pauseSignalCh <- struct{}{}
				}
			}
		}
	}
}

func (g *Gofermart) worker(ctx context.Context, wgWait *sync.WaitGroup, workerID int, resultCh chan<- int, pauseSignalCh chan struct{}) {

	defer wgWait.Done()

	for {
		select {
		case <-ctx.Done():
			g.loggerer.Debugw("stop worker by context",
				"workerId", workerID,
			)
			return
		case <-pauseSignalCh:
			g.loggerer.Debugw("paused worker number",
				"workerID", workerID,
			)
			time.Sleep(SleepCheckTime)
		case order := <-g.inputCh:
			g.loggerer.Debugw("start processing order",
				"workerId", workerID,
				"order number", order.OrderNumber,
				"order status", order.OrderStatus,
			)
			ctx1, cancel1 := context.WithTimeout(ctx, time.Second)
			res, status := g.orderChecker.OrderCheck(ctx1, order.OrderNumber)
			cancel1()
			g.loggerer.Debugw("return order from OrderCheck",
				"workerId", workerID,
				"order number", res.OrderNumber,
				"order status", res.OrderStatus,
				"status", status,
			)

			if status != http.StatusOK {
				g.wg.Add(1)
				go g.addToInputCh(order, RepetitiveCheckTime)
				resultCh <- status
				break
			}

			if res.OrderStatus == storage.OrderStatusNew || res.OrderStatus == storage.OrderStatusProcessing {
				g.wg.Add(1)
				go g.addToInputCh(order, RepetitiveCheckTime)
				resultCh <- status
				break
			}

			if res.OrderStatus != order.OrderStatus {
				ctx2, cancel2 := context.WithTimeout(ctx, time.Second)
				err := g.storager.UpdateOrder(ctx2, order.UserLogin, order.OrderNumber, res.OrderStatus, res.BonusPoints)
				cancel2()
				if err != nil {
					g.loggerer.Debugw("can't update order status in database",
						"workerId", workerID,
						"order number", res.OrderNumber,
						"order status", res.OrderStatus,
						"error", err,
					)
					g.wg.Add(1)
					go g.addToInputCh(order, RepetitiveCheckTime)
					resultCh <- status
					break
				}


				resultCh <- status
				g.loggerer.Debugw("finished processing order successful",
					"workerId", workerID,
					"order number", res.OrderNumber,
					"order status", res.OrderStatus,
					"accrual", res.BonusPoints,
					"status", status,
				)
			}
		}
	}
}

// Add order to input channel
func (g *Gofermart) addToInputCh(order storage.Order, waitTime time.Duration) {

	time.Sleep(waitTime)
	count++

	defer func() {
		g.wg.Done()
		count--
	}()

	select {
	case <-g.doneCh:
		return
	case g.inputCh <- order:
		g.loggerer.Debugw("add order to inputCh",
			"order number", order.OrderNumber,
		)
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

	var res int
	const k = 9
	dataLenght := len(data)
	startNum := dataLenght % 2

	if startNum != 0 {
		res += data[0]
	}

	for ; startNum < dataLenght-1; startNum += 2 {
		data[startNum] *= 2
		if data[startNum] >= k {
			data[startNum] -= k
		}
		res += data[startNum] + data[startNum+1]
	}

	return res%10 == 0
}

//	5 7 4 7 0 7 2 1 3 5 8
//	5 5 4 5 0 5 2 2 3 1
