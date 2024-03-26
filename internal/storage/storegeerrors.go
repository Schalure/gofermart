package storage

import "errors"

var (
	ErrLoginAlreadyExists       = errors.New("login already exists")
	ErrOrderNumberAlreadyExists = errors.New("order number already exists")

	ErrInsufficientFunds = errors.New("insufficient funds")
)
