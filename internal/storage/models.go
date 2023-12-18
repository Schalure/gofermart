package storage

import (
	"encoding/json"
	"time"
)

type User struct {
	Login           string  `json:"login" db:"login"`
	Password        string  `json:"password" db:"password"`
	LoyaltyPoints   float64 `json:"loyalty_points" db:"loyalty_points"`
	WithdrawnPoints float64 `json:"withdrawn_points" db:"withdrawn_points"`
}

func (u User) String() string {

	data, err := json.Marshal(&u)
	if err != nil {
		return ""
	}
	return string(data)
}

// Доступные статусы обработки расчётов:
type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"        //	заказ загружен в систему, но не попал в обработку
	OrderStatusProcessing             = "PROCESSING" //	вознаграждение за заказ рассчитывается
	OrderStatusInvalid                = "INVALID"    //	система расчёта вознаграждений отказала в расчёте
	OrderStatusProcessed              = "PROCESSED"  //	данные по заказу проверены и информация о расчёте успешно получена
)

type Order struct {
	OrderNumber string      `json:"number" db:"order_number"`
	OrderStatus OrderStatus `json:"status" db:"order_status"`
	BonusPoints float64     `json:"accural" db:"bonus_points"`
	UploadedAt  time.Time   `json:"uploaded_at" db:"uploaded_at"`
	UserLogin   string      `json:"login" db:"user_login"`
}

func (o Order) String() string {

	data, err := json.Marshal(&o)
	if err != nil {
		return ""
	}
	return string(data)
}
