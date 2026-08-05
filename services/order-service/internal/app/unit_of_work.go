package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/order-service/internal/domain"
	"github.com/google/uuid"
)

type OrderStore interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrderByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	ListOrdersByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error
	SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload []byte) error
}

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(store OrderStore) error) error
}
