package domain

import (
	"time"

	"github.com/google/uuid"
)

type InventoryItem struct {
	AggregateRoot
	ProductID         uuid.UUID
	AvailableQuantity int32
	ReservedQuantity  int32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type StockReservedEvent struct {
	ReservationID string    `json:"reservation_id"`
	ProductID     string    `json:"product_id"`
	Quantity      int32     `json:"quantity"`
	Timestamp     time.Time `json:"timestamp"`
}

func (e StockReservedEvent) EventType() string     { return "StockReserved" }
func (e StockReservedEvent) OccurredAt() time.Time { return e.Timestamp }

type StockReleasedEvent struct {
	ReservationID string    `json:"reservation_id"`
	ProductID     string    `json:"product_id"`
	Quantity      int32     `json:"quantity"`
	Timestamp     time.Time `json:"timestamp"`
}

func (e StockReleasedEvent) EventType() string     { return "StockReleased" }
func (e StockReleasedEvent) OccurredAt() time.Time { return e.Timestamp }

type StockReservation struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewInventoryItem(productID uuid.UUID, initialQuantity int32) (*InventoryItem, error) {
	if initialQuantity < 0 {
		return nil, ErrInvalidQuantity
	}
	now := time.Now().UTC()
	return &InventoryItem{
		ProductID:         productID,
		AvailableQuantity: initialQuantity,
		ReservedQuantity:  0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

func (i *InventoryItem) Reserve(quantity int32) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if i.AvailableQuantity < quantity {
		return ErrInsufficientStock
	}

	i.AvailableQuantity -= quantity
	i.ReservedQuantity += quantity
	i.UpdatedAt = time.Now().UTC()
	return nil
}

func (i *InventoryItem) Release(quantity int32) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if i.ReservedQuantity < quantity {
		i.AvailableQuantity += i.ReservedQuantity
		i.ReservedQuantity = 0
	} else {
		i.ReservedQuantity -= quantity
		i.AvailableQuantity += quantity
	}
	i.UpdatedAt = time.Now().UTC()
	return nil
}
