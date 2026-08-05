package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusConfirmed OrderStatus = "CONFIRMED"
	StatusCancelled OrderStatus = "CANCELLED"
	StatusPaid      OrderStatus = "PAID"
)

type OrderItem struct {
	ID         uuid.UUID
	OrderID    uuid.UUID
	ProductID  uuid.UUID
	Quantity   int32
	PriceCents int64
	CreatedAt  time.Time
}

type Order struct {
	AggregateRoot
	ID         uuid.UUID
	UserID     uuid.UUID
	Items      []*OrderItem
	TotalCents int64
	Status     OrderStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type OrderCreatedEvent struct {
	OrderID    string    `json:"order_id"`
	UserID     string    `json:"user_id"`
	TotalCents int64     `json:"total_cents"`
	Timestamp  time.Time `json:"timestamp"`
}

func (e OrderCreatedEvent) EventType() string     { return "OrderCreated" }
func (e OrderCreatedEvent) OccurredAt() time.Time { return e.Timestamp }

type OrderCancelledEvent struct {
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

func (e OrderCancelledEvent) EventType() string     { return "OrderCancelled" }
func (e OrderCancelledEvent) OccurredAt() time.Time { return e.Timestamp }

type OrderPaidEvent struct {
	OrderID      string    `json:"order_id"`
	UserID       string    `json:"user_id"`
	TotalCents   int64     `json:"total_cents"`
	TransactionID string   `json:"transaction_id"`
	Timestamp    time.Time `json:"timestamp"`
}

func (e OrderPaidEvent) EventType() string     { return "OrderPaid" }
func (e OrderPaidEvent) OccurredAt() time.Time { return e.Timestamp }

func NewOrder(userID uuid.UUID, itemInputs []struct {
	ProductID  uuid.UUID
	Quantity   int32
	PriceCents int64
}) (*Order, error) {
	if len(itemInputs) == 0 {
		return nil, ErrEmptyOrder
	}

	orderID := uuid.New()
	now := time.Now().UTC()
	var totalCents int64
	items := make([]*OrderItem, 0, len(itemInputs))

	for _, input := range itemInputs {
		if input.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
		if input.PriceCents < 0 {
			return nil, ErrInvalidPrice
		}

		itemTotal := int64(input.Quantity) * input.PriceCents
		totalCents += itemTotal

		items = append(items, &OrderItem{
			ID:         uuid.New(),
			OrderID:    orderID,
			ProductID:  input.ProductID,
			Quantity:   input.Quantity,
			PriceCents: input.PriceCents,
			CreatedAt:  now,
		})
	}

	order := &Order{
		ID:         orderID,
		UserID:     userID,
		Items:      items,
		TotalCents: totalCents,
		Status:     StatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	order.RecordEvent(OrderCreatedEvent{
		OrderID:    orderID.String(),
		UserID:     userID.String(),
		TotalCents: totalCents,
		Timestamp:  now,
	})

	return order, nil
}

func (o *Order) Cancel(reason string) error {
	if o.Status == StatusCancelled || o.Status == StatusPaid {
		return ErrCannotCancelOrder
	}

	now := time.Now().UTC()
	o.Status = StatusCancelled
	o.UpdatedAt = now

	o.RecordEvent(OrderCancelledEvent{
		OrderID:   o.ID.String(),
		UserID:    o.UserID.String(),
		Reason:    reason,
		Timestamp: now,
	})

	return nil
}

func (o *Order) MarkAsPaid(transactionID string) error {
	if o.Status == StatusPaid {
		return nil
	}
	if o.Status == StatusCancelled {
		return ErrCannotMarkPaidCancelled
	}

	now := time.Now().UTC()
	o.Status = StatusPaid
	o.UpdatedAt = now

	o.RecordEvent(OrderPaidEvent{
		OrderID:       o.ID.String(),
		UserID:        o.UserID.String(),
		TotalCents:    o.TotalCents,
		TransactionID: transactionID,
		Timestamp:     now,
	})

	return nil
}
