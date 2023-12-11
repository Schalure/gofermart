package gofermart

import (
	"fmt"

	"github.com/Schalure/gofermart/internal/storage"
)

func (g *Gofermart) LoadOrder(login, orderNumber string) error {



	return nil
}

func (g *Gofermart) GetOrders(login string) ([]storage.Order, error) {


	return nil, nil
}

func (g *Gofermart) isOrderValid(orderNumber string) bool {

	if !g.validOrderNumber.MatchString(orderNumber) {
		return false
	}

	var data []int
	for _, b := range orderNumber {
		data = append(data, int(b))
	}
	fmt.Println(data)

	
	return true
}
