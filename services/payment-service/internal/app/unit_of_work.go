package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/payment-service/internal/domain"
	"github.com/google/uuid"
)

type PaymentStore interface {
	CreatePayment(ctx context.Context, payment *domain.Payment) error
	GetPaymentByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error)
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus) error
	SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload []byte) error
}

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(store PaymentStore) error) error
}
