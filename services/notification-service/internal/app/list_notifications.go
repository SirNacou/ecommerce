package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/notification-service/internal/domain"
	"github.com/google/uuid"
)

type ListNotificationsQueryHandler struct {
	uow UnitOfWork
}

func NewListNotificationsQueryHandler(uow UnitOfWork) *ListNotificationsQueryHandler {
	return &ListNotificationsQueryHandler{uow: uow}
}

func (h *ListNotificationsQueryHandler) Handle(ctx context.Context, userIDStr string, limit, offset int32) ([]*domain.Notification, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, domain.ErrNotificationNotFound
	}

	var notifications []*domain.Notification
	err = h.uow.Execute(ctx, func(store NotificationStore) error {
		var err error
		notifications, err = store.ListNotificationsByUserID(ctx, userID, limit, offset)
		return err
	})
	if err != nil {
		return nil, err
	}

	return notifications, nil
}
