package gofermaterrors

import "errors"

var (
	ErrInternal = errors.New("something went wrong")
	ErrNoData   = errors.New("no data")

	//	Users errors
	ErrLoginAlreadyTaken    = errors.New("the login is already taken")
	ErrInvalidLogin         = errors.New("invalid login")
	ErrPasswordShort        = errors.New("the password is too short, make up a password of 8 characters or more")
	ErrPasswordBad          = errors.New("the password can only be made of characters: 0-9, a-z, A-Z")
	ErrInvalidLoginPassword = errors.New("invalid login or password")

	//	orders errors
	ErrInvalidOrderNumber         = errors.New("invalid order number")
	ErrDublicateOrderNumberByUser = errors.New("the order number has already been uploaded by this user")
	ErrDublicateOrderNumber       = errors.New("the order number has already been uploaded by another user")
	ErrNoOrdersForThisUser        = errors.New("no orders for this user")
	ErrInsufficientFunds          = errors.New("insufficient funds")
	ErrNoOrdersForPoints          = errors.New("no orders for points")
)
