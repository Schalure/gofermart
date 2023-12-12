package gofermart

import (
	"strconv"
	"strings"

	"github.com/Schalure/gofermart/internal/gofermart/gofermaterrors"
	"github.com/Schalure/gofermart/internal/storage"
)

func (g *Gofermart) LoadOrder(login, orderNumber string) error {

	if !g.isOrderValid(orderNumber) {
		return gofermaterrors.InvalidOrderNumber
	}

	return nil
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
