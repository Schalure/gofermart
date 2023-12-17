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

func (s *LoyaltySystem) OrderCheck(ctx context.Context, ordernumber string) (storage.Order, int) {

	type result struct {
		Order   string      `json:"order"`
		Status  OrderStatus `json:"status"`
		Accrual float64         `json:"accrual"`
	}
	var res result

	resp, err := s.client.R().
		SetContext(ctx).
		SetResult(&res).
		SetPathParams(map[string]string{"number": ordernumber}).
		Get(s.queryString)

	if err != nil {
		return storage.Order{}, 0
	}

	return storage.Order{
		OrderNumber: res.Order,
		OrderStatus: storage.OrderStatus(res.Status),
		BonusPoints: res.Accrual,
	}, resp.StatusCode()
}

