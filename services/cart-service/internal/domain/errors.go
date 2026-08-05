package domain

import "errors"

var (
	ErrCartNotFound    = errors.New("cart not found")
	ErrItemNotFound    = errors.New("item not found in cart")
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")
	ErrInvalidPrice    = errors.New("price cannot be negative")
	ErrProductNotFound = errors.New("product not found")
)
