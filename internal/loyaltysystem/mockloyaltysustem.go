package loyaltysystem

import (
	"context"
	"math/rand"

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

	statusCode := 200


	if statusCode == 200 {
		statusBig := rand.Intn(2)

		switch statusBig {
		case 0: res.Status = "PROCESSING"
		case 1: res.Status = "PROCESSED"
		default: panic("OrderCheck")
		}

		if res.Status == "PROCESSED" {
			bonusPoints := 100 + rand.Float64() * (500 - 100)
			res.Accrual = bonusPoints
		}
	
	}
	
	return storage.Order{
		OrderNumber: ordernumber,
		OrderStatus: storage.OrderStatus(res.Status),
		BonusPoints: res.Accrual,
	}, statusCode

}