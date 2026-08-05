package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	CategoryID  uuid.UUID `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PriceCents  int64     `json:"price_cents"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewProduct(categoryID uuid.UUID, name, description string, priceCents int64) (*Product, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	if priceCents <= 0 {
		return nil, ErrInvalidPrice
	}

	now := time.Now().UTC()
	return &Product{
		ID:          uuid.New(),
		CategoryID:  categoryID,
		Name:        name,
		Description: description,
		PriceCents:  priceCents,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
