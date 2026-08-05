package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/notification-service/internal/domain"
	"github.com/google/uuid"
)

type NotificationStore interface {
	CreateNotification(ctx context.Context, n *domain.Notification) error
	GetNotificationByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	UpdateNotificationStatus(ctx context.Context, id uuid.UUID, status domain.NotificationStatus) error
	SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload []byte) error
}

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(store NotificationStore) error) error
}
