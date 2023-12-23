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
		Accrual float64     `json:"accrual"`
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

	// var statusCode int
	// statusCodeBig, err := rand.Int(rand.Reader, big.NewInt(3))
	// if err != nil {
	// 	panic("OrderCheck")
	// }
	// switch statusCodeBig.Int64() {
	// case 0: statusCode = 200
	// case 1: statusCode = 204
	// case 2: statusCode = 429
	// case 3: statusCode = 500
	// default: panic("OrderCheck")
	// }

	// if statusCode == 200 {
	// 	statusBig, err := rand.Int(rand.Reader, big.NewInt(3))
	// 	if err != nil {
	// 		panic("OrderCheck")
	// 	}
	// 	switch statusBig.Int64() {
	// 	case 0: res.Status = "NEW"
	// 	case 1, 3: res.Status = "PROCESSED"
	// 	case 2: res.Status = "PROCESSING"
	// 	default: panic("OrderCheck")
	// 	}

	// 	if res.Status == "PROCESSED" {
	// 		bonusPointsBig, err := rand.Int(rand.Reader, big.NewInt(500))
	// 		if err != nil {
	// 			panic("OrderCheck")
	// 		}
	// 		res.Accrual = float64(bonusPointsBig.Int64())
	// 	}
	
	// }

	// return storage.Order{
	// 	OrderNumber: ordernumber,
	// 	OrderStatus: storage.OrderStatus(res.Status),
	// 	BonusPoints: res.Accrual,
	// }, statusCode

}
