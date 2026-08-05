package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "PENDING"
	StatusCompleted PaymentStatus = "COMPLETED"
	StatusFailed    PaymentStatus = "FAILED"
	StatusRefunded  PaymentStatus = "REFUNDED"
)

type Payment struct {
	AggregateRoot
	ID            uuid.UUID
	OrderID       uuid.UUID
	UserID        uuid.UUID
	AmountCents   int64
	Currency      string
	Status        PaymentStatus
	PaymentMethod string
	TransactionID string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PaymentProcessedEvent struct {
	PaymentID     string    `json:"payment_id"`
	OrderID       string    `json:"order_id"`
	UserID        string    `json:"user_id"`
	AmountCents   int64     `json:"amount_cents"`
	Status        string    `json:"status"`
	TransactionID string    `json:"transaction_id"`
	Timestamp     time.Time `json:"timestamp"`
}

func (e PaymentProcessedEvent) EventType() string     { return "PaymentProcessed" }
func (e PaymentProcessedEvent) OccurredAt() time.Time { return e.Timestamp }

type PaymentRefundedEvent struct {
	PaymentID string    `json:"payment_id"`
	OrderID   string    `json:"order_id"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

func (e PaymentRefundedEvent) EventType() string     { return "PaymentRefunded" }
func (e PaymentRefundedEvent) OccurredAt() time.Time { return e.Timestamp }

func NewPayment(
	orderID, userID uuid.UUID,
	amountCents int64,
	currency, method string,
) (*Payment, error) {
	if amountCents <= 0 {
		return nil, ErrInvalidAmount
	}

	if currency == "" {
		currency = "USD"
	}

	paymentID := uuid.New()
	transactionID := "txn_" + uuid.New().String()[:18]
	now := time.Now().UTC()

	payment := &Payment{
		ID:            paymentID,
		OrderID:       orderID,
		UserID:        userID,
		AmountCents:   amountCents,
		Currency:      currency,
		Status:        StatusCompleted,
		PaymentMethod: method,
		TransactionID: transactionID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	payment.RecordEvent(PaymentProcessedEvent{
		PaymentID:     paymentID.String(),
		OrderID:       orderID.String(),
		UserID:        userID.String(),
		AmountCents:   amountCents,
		Status:        string(StatusCompleted),
		TransactionID: transactionID,
		Timestamp:     now,
	})

	return payment, nil
}

func (p *Payment) Refund(reason string) error {
	if p.Status != StatusCompleted {
		return ErrCannotRefund
	}

	now := time.Now().UTC()
	p.Status = StatusRefunded
	p.UpdatedAt = now

	p.RecordEvent(PaymentRefundedEvent{
		PaymentID: p.ID.String(),
		OrderID:   p.OrderID.String(),
		Reason:    reason,
		Timestamp: now,
	})

	return nil
}
