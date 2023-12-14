package gofermart

import (
	"fmt"

	"github.com/Schalure/gofermart/internal/storage"
)

const ordersCashSize = 256

//	Cash method. Add order to cash. Return error if cash is full
func (g *Gofermart) addOrderToCash(order storage.Order) error {

	if len(g.orderCash) < ordersCashSize {
		g.orderCash[order.OrderNumber] = order
		return nil
	}
	return fmt.Errorf("order cash if full")
}

//	Cash method. Get order from cash. Return false if order not found in cash
func (g *Gofermart) getOrderFromCash(orderNumber string) (storage.Order, bool) {
	order, ok := g.orderCash[orderNumber]
	return order, ok
}

//	Cash method. Delete order from cash. Return false if order not found in cash
func (g *Gofermart) deleteOrderFromCash(orderNumber string) bool {

	_, ok := g.orderCash[orderNumber]
	if ok {
		delete(g.orderCash, orderNumber)
	}
	return  ok
}

//	Cash method. Return true if order cash is not full, and count of free cells
func (g *Gofermart) isCashNotFull() (int, bool) {

	l := len(g.orderCash)
	if l == ordersCashSize {
		return 0, false
	}
	return ordersCashSize - l, true
}