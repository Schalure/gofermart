package storage

import (
	"encoding/json"

	"github.com/jackc/pgx/pgtype"
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
	OrderStatusProcessing OrderStatus             = "PROCESSING" //	вознаграждение за заказ рассчитывается
	OrderStatusInvalid OrderStatus                = "INVALID"    //	система расчёта вознаграждений отказала в расчёте
	OrderStatusProcessed OrderStatus              = "PROCESSED"  //	данные по заказу проверены и информация о расчёте успешно получена
)

type Order struct {
	OrderNumber string      `json:"number" db:"order_number"`
	OrderStatus OrderStatus `json:"status" db:"order_status"`
//	UploadedOrder time.Time `json:"uploaded_order" db:"uploaded_order"`
	UploadedOrder pgtype.Timestamptz `json:"uploaded_order" db:"uploaded_order"`
	BonusPoints float64     `json:"accural" db:"bonus_points"`
//	UploadedBonus time.Time `json:"uploaded_bonus" db:"uploaded_bonus"`
	UploadedBonus pgtype.Timestamptz `json:"uploaded_bonus" db:"uploaded_bonus"`
	UserLogin   string      `json:"login" db:"user_login"`
}

func (o Order) String() string {

	data, err := json.Marshal(&o)
	if err != nil {
		return ""
	}
	return string(data)
}
