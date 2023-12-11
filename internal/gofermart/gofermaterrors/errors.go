package gofermaterrors

import "errors"

var (
	Internal             = errors.New("something went wrong")
	LoginAlreadyTaken    = errors.New("the login is already taken")
	InvalidLogin         = errors.New("invalid login")
	PasswordShort        = errors.New("the password is too short, make up a password of 8 characters or more")
	PasswordBad          = errors.New("the password can only be made of characters: 0-9, a-z, A-Z")
	InvalidLoginPassword = errors.New("invalid login or password")
)
