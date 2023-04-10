package store

import "errors"

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrUserAlreadyExists = errors.New("user with email already exists")
	ErrInsufficientFunds = errors.New("not enough funds on source account")
)
