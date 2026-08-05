package domain

import "errors"

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrInvalidName       = errors.New("name cannot be empty")
	ErrInvalidPrice      = errors.New("price must be greater than zero")
	ErrNegativeStock     = errors.New("stock cannot be negative")
	ErrItemNotFound      = errors.New("inventory item not found")
	ErrInsufficientStock = errors.New("insufficient available stock")
	ErrInvalidQuantity   = errors.New("quantity must be greater than zero")
)
