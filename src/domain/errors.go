package domain

import (
	"errors"
)

// Ошибки для аккаунтов
var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrAccountIncorrect  = errors.New("account incorrect data")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidType       = errors.New("invalid account type")
	ErrReservationFailed = errors.New("reservation failed")
	ErrReleaseFailed     = errors.New("release failed")
	ErrNotEnoughFunds    = errors.New("not Enough Funds")
)

// Ошибки для платежей
var (
	ErrPaymentIncorrectStatus = errors.New("payment incorrect status")
	ErrPaymentNotFound        = errors.New("payment not found")
	ErrAlreadyExists          = errors.New("payment already exists")
	ErrInvalidStatus          = errors.New("invalid payment status")
	ErrNotEnougthFunds        = errors.New("not enought funds")
	ErrInvalidAmount          = errors.New("invalid Amount")
	ErrSameAccount            = errors.New("same Account")
)
