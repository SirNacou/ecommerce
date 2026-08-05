package domain

import "errors"

var (
	ErrOrderNotFound         = errors.New("order not found")
	ErrProductNotFound       = errors.New("product not found")
	ErrEmptyOrder            = errors.New("order must contain at least one item")
	ErrInvalidQuantity       = errors.New("quantity must be greater than zero")
	ErrInvalidPrice          = errors.New("price cannot be negative")
	ErrCannotCancelOrder     = errors.New("order cannot be cancelled in its current state")
	ErrCannotMarkPaidCancelled = errors.New("cancelled order cannot be marked as paid")
)
