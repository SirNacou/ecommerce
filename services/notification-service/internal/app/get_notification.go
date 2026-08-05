package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/notification-service/internal/domain"
	"github.com/google/uuid"
)

type GetNotificationQueryHandler struct {
	uow UnitOfWork
}

func NewGetNotificationQueryHandler(uow UnitOfWork) *GetNotificationQueryHandler {
	return &GetNotificationQueryHandler{uow: uow}
}

func (h *GetNotificationQueryHandler) Handle(ctx context.Context, idStr, userIDStr string) (*domain.Notification, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, domain.ErrNotificationNotFound
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, domain.ErrNotificationNotFound
	}

	var notification *domain.Notification
	err = h.uow.Execute(ctx, func(store NotificationStore) error {
		var err error
		notification, err = store.GetNotificationByID(ctx, id)
		return err
	})
	if err != nil {
		return nil, err
	}

	if notification.UserID != userID {
		return nil, domain.ErrNotificationNotFound
	}

	return notification, nil
}
