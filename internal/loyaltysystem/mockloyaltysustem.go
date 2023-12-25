package loyaltysystem

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/Schalure/gofermart/internal/storage"
)

type MockLoyaltySystem struct {

}

func NewMockLoyaltySystem() *MockLoyaltySystem {

	return &MockLoyaltySystem{}
}

func (s *MockLoyaltySystem) OrderCheck(ctx context.Context, ordernumber string) (storage.Order, int) {

	type result struct {
		Order   string      `json:"order"`
		Status  OrderStatus `json:"status"`
		Accrual float64     `json:"accrual"`
	}
	var res result

	var statusCode int
	statusCodeBig, err := rand.Int(rand.Reader, big.NewInt(3))
	if err != nil {
		panic("OrderCheck")
	}
	switch statusCodeBig.Int64() {
	case 0: statusCode = 200
	case 1: statusCode = 204
	case 2: statusCode = 429
	case 3: statusCode = 500
	default: panic("OrderCheck")
	}

	if statusCode == 200 {
		statusBig, err := rand.Int(rand.Reader, big.NewInt(3))
		if err != nil {
			panic("OrderCheck")
		}
		switch statusBig.Int64() {
		case 0: res.Status = "NEW"
		case 1, 3: res.Status = "PROCESSED"
		case 2: res.Status = "PROCESSING"
		default: panic("OrderCheck")
		}

		if res.Status == "PROCESSED" {
			bonusPointsBig, err := rand.Int(rand.Reader, big.NewInt(500))
			if err != nil {
				panic("OrderCheck")
			}
			res.Accrual = float64(bonusPointsBig.Int64())
		}
	
	}
	
	return storage.Order{
		OrderNumber: ordernumber,
		OrderStatus: storage.OrderStatus(res.Status),
		BonusPoints: res.Accrual,
	}, statusCode

}