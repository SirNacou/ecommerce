package app

import (
	"context"
	"encoding/json"

	"github.com/SirNacou/ecommerce/services/notification-service/internal/domain"
	"github.com/google/uuid"
)

type SendNotificationCommand struct {
	UserID    string
	Channel   domain.Channel
	Recipient string
	Subject   string
	Body      string
}

type SendNotificationCommandHandler struct {
	uow UnitOfWork
}

func NewSendNotificationCommandHandler(uow UnitOfWork) *SendNotificationCommandHandler {
	return &SendNotificationCommandHandler{uow: uow}
}

func (h *SendNotificationCommandHandler) Handle(ctx context.Context, cmd SendNotificationCommand) (*domain.Notification, error) {
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, domain.ErrNotificationNotFound
	}

	notification, err := domain.NewNotification(userID, cmd.Channel, cmd.Recipient, cmd.Subject, cmd.Body)
	if err != nil {
		return nil, err
	}

	err = h.uow.Execute(ctx, func(store NotificationStore) error {
		if err := store.CreateNotification(ctx, notification); err != nil {
			return err
		}

		for _, event := range notification.PopEvents() {
			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}
			if err := store.SaveOutboxEvent(ctx, "Notification", notification.ID.String(), event.EventType(), payload); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return notification, nil
}