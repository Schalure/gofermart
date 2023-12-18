package gofermaterrors

import "errors"

var (
	Internal = errors.New("something went wrong")
	NoData   = errors.New("no data")

	//	Users errors
	LoginAlreadyTaken    = errors.New("the login is already taken")
	InvalidLogin         = errors.New("invalid login")
	PasswordShort        = errors.New("the password is too short, make up a password of 8 characters or more")
	PasswordBad          = errors.New("the password can only be made of characters: 0-9, a-z, A-Z")
	InvalidLoginPassword = errors.New("invalid login or password")

	//	orders errors
	InvalidOrderNumber         = errors.New("invalid order number")
	DublicateOrderNumberByUser = errors.New("the order number has already been uploaded by this user")
	DublicateOrderNumber       = errors.New("the order number has already been uploaded by another user")
	NoOrdersForThisUser        = errors.New("no orders for this user")
	InsufficientFunds          = errors.New("insufficient funds")
	NoOrdersForPoints          = errors.New("no orders for points")
)
