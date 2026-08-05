package domain

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID         uuid.UUID `json:"id"`
	CartID     uuid.UUID `json:"cart_id"`
	ProductID  uuid.UUID `json:"product_id"`
	Quantity   int32     `json:"quantity"`
	PriceCents int64     `json:"price_cents"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Cart struct {
	ID        uuid.UUID   `json:"id"`
	UserID    uuid.UUID   `json:"user_id"`
	Items     []*CartItem `json:"items"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func NewCart(userID uuid.UUID) *Cart {
	now := time.Now().UTC()
	return &Cart{
		ID:        uuid.New(),
		UserID:    userID,
		Items:     make([]*CartItem, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (c *Cart) TotalCents() int64 {
	var total int64
	for _, item := range c.Items {
		total += int64(item.Quantity) * item.PriceCents
	}
	return total
}
