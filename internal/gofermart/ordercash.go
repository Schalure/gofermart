package gofermart

import (
	"sync"

	"github.com/Schalure/gofermart/internal/storage"
)

type OrderCash struct {

	cash map[string]storage.Order

	mu sync.RWMutex
}

func NewOrderCash() *OrderCash {

	return &OrderCash{
		cash: make(map[string]storage.Order),
	}
}

func (c *OrderCash) add(order storage.Order) {
	c.mu.Lock()
	c.cash[order.OrderNumber] = order
	c.mu.Unlock()
}

func (c *OrderCash) get() storage.Order {
	var order storage.Order
	c.mu.RLock()
	for _, order = range c.cash{
		break
	}
	c.mu.RUnlock()
	return order
}

func (c *OrderCash) delete(orderNumber string) {
	c.mu.Lock()
	delete(c.cash, orderNumber)
	c.mu.Unlock()
}

func (c *OrderCash) isEmpty() bool {
	c.mu.RLock()
	l := len(c.cash)
	c.mu.RUnlock()
	return l == 0
}