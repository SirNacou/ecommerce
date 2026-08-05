package domain

import "errors"

var (
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrInvalidAmount    = errors.New("payment amount must be greater than zero")
	ErrAlreadyProcessed = errors.New("payment has already been processed for this order")
	ErrCannotRefund     = errors.New("only completed payments can be refunded")
	ErrPaymentFailed    = errors.New("payment authorization failed")
)
