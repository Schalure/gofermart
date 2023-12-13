package loyaltysystem

import (
	"context"

	"github.com/Schalure/gofermart/internal/storage"
	"github.com/go-resty/resty/v2"
)

const orderNumberParName = "number"

// Order status constants
type OrderStatus string

const (
	Registered OrderStatus = "REGISTERED" //	заказ зарегистрирован, но не начисление не рассчитано;
	Invalid    OrderStatus = "INVALID"    //	заказ не принят к расчёту, и вознаграждение не будет начислено;
	Processing OrderStatus = "PROCESSING" //	расчёт начисления в процессе;
	Processed  OrderStatus = "PROCESSED"  //	расчёт начисления окончен;
)

// LoyaltySystem object struct
type LoyaltySystem struct {
	client      *resty.Client
	host        string
	queryString string
}

// Constructor of LoyaltySystem object struct
func NewLoyaltySystem(host string) *LoyaltySystem {

	queryString := host + "/api/orders/{" + orderNumberParName + "}"

	return &LoyaltySystem{
		client:      resty.New(),
		host:        host,
		queryString: queryString,
	}
}

func (s *LoyaltySystem) OrderCheck(ctx context.Context, order storage.Order, resultCh chan<- storage.Order) {

	type result struct {
		Order   string      `json:"order"`
		Status  OrderStatus `json:"status"`
		Accrual int         `json:"accrual"`
	}
	var res result

	_, err := s.client.R().
		SetContext(ctx).
		SetResult(&res).
		SetPathParams(map[string]string{orderNumberParName: order.OrderNumber}).
		Get(s.queryString)

	if err != nil {
		return
	}

	resultCh <- storage.Order{
		OrderNumber: res.Order,
		OrderStatus: storage.OrderStatus(res.Status),
		BonusPoints: res.Accrual,
		UserLogin: order.UserLogin,
		UploadedAt: order.UploadedAt,
	}
}
