package storage

import (
	"encoding/json"
	"time"
)

type User struct {
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}

func (u User) String() string {

	data, err := json.Marshal(&u)
	if err != nil {
		return ""
	}
	return string(data)
}

type Order struct {
	OrderNumber string `json:"number" db:"order_number"`
	OrderStatus string `json:"status" db:"order_status"`
	BonusPoints string `json:"accural" db:"bonus_points"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
	UserLogin string `json:"login" db:"user_login"`
}

func (o Order) String() string {

	data, err := json.Marshal(&o)
	if err != nil {
		return ""
	}
	return string(data)
}