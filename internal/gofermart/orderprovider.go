package gofermart

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"time"

// 	"github.com/Schalure/gofermart/internal/storage"
// )

// //	Order provider constants
// const (
// 	orderProviderCashSize = 256
// 	countOfJobs = 20
// )

// //	Order provider object struct
// type OrderProvider struct {
// 	orderCash map[string]storage.Order
// 	OrderStatusChecker
// }

// //	Constructor of OrderProvider
// func newOrderProvider(statusChecker OrderStatusChecker) *OrderProvider {

// 	return &OrderProvider{
// 		orderCash: make(map[string]storage.Order, orderProviderCashSize),
// 		OrderStatusChecker: statusChecker,
// 	}
// }

// //	Start order provider task
// func (p *OrderProvider) Run(ctx context.Context) {

// }

// //	Worker читает базу данных
// func (p *OrderProvider) worker(ctx context.Context) {

// 	jobsCh := make(chan storage.Order, countOfJobs)
// 	resultCh := make(chan storage.Order, countOfJobs)
// 	var wg sync.WaitGroup

// 	for i := 0; i < countOfJobs; i++ {

// 		wg.Add(1)
// 		go func (jobsCh <-chan storage.Order, resultCh chan<- storage.Order) {
// 			defer wg.Done()
// 			for order := range jobsCh {
// 				ctx, cancel := context.WithTimeout(ctx, time.Second * 1)
// 				resultOrder, err := p.OrderStatusCheck(ctx, order.OrderNumber)
// 				cancel()
// 				if err == nil {
// 					order.BonusPoints = resultOrder.BonusPoints
// 					order.OrderStatus = resultOrder.OrderStatus
// 					resultCh <- order
// 				}
// 			}
// 		}(jobsCh, resultCh)
// 	}

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			close(jobsCh)
// 			wg.Wait()
// 			close(resultCh)
// 			return
// 		case order := <-resultCh:

// 		}
// 	}
// }

// //	Cash method. Add order to cash. Return error if cash is full
// func (p *OrderProvider) addOrderToCash(order storage.Order) error {

// 	if len(p.orderCash) < orderProviderCashSize {
// 		p.orderCash[order.OrderNumber] = order
// 		return nil
// 	}
// 	return fmt.Errorf("order cash if full")
// }

// //	Cash method. Get order from cash. Return false if order not found in cash
// func (p *OrderProvider) getOrderFromCash(orderNumber string) (storage.Order, bool) {
// 	order, ok := p.orderCash[orderNumber]
// 	return order, ok
// }

// //	Cash method. Delete order from cash. Return false if order not found in cash
// func (p *OrderProvider) deleteOrderFromCash(orderNumber string) bool {

// 	_, ok := p.orderCash[orderNumber]
// 	if ok {
// 		delete(p.orderCash, orderNumber)
// 	}
// 	return  ok
// }

// //	Cash method. Return true if order cash is not full, and count of free cells
// func (p *OrderProvider) isCashNotFull() (int, bool) {

// 	l := len(p.orderCash)
// 	if l == orderProviderCashSize {
// 		return 0, false
// 	}
// 	return orderProviderCashSize - l, true
// }

