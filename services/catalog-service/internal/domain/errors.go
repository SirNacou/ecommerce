package domain

import "errors"

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrInvalidName      = errors.New("name cannot be empty")
	ErrInvalidPrice     = errors.New("price must be greater than zero")
	ErrNegativeStock    = errors.New("stock cannot be negative")
)
