package app

import (
	"context"
	"encoding/json"

	"github.com/SirNacou/ecommerce/services/order-service/internal/domain"
	"github.com/google/uuid"
)

type CancelOrderCommand struct {
	OrderID string
	UserID  string
	Reason  string
}

type CancelOrderCommandHandler struct {
	uow UnitOfWork
}

func NewCancelOrderCommandHandler(uow UnitOfWork) *CancelOrderCommandHandler {
	return &CancelOrderCommandHandler{uow: uow}
}

func (h *CancelOrderCommandHandler) Handle(ctx context.Context, cmd CancelOrderCommand) (*domain.Order, error) {
	orderID, err := uuid.Parse(cmd.OrderID)
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}

	var order *domain.Order
	err = h.uow.Execute(ctx, func(store OrderStore) error {
		var err error
		order, err = store.GetOrderByID(ctx, orderID)
		if err != nil {
			return err
		}

		if order.UserID != userID {
			return domain.ErrOrderNotFound
		}

		if err := order.Cancel(cmd.Reason); err != nil {
			return err
		}

		if err := store.UpdateOrderStatus(ctx, order.ID, order.Status); err != nil {
			return err
		}

		for _, event := range order.PopEvents() {
			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}
			if err := store.SaveOutboxEvent(ctx, "Order", order.ID.String(), event.EventType(), payload); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}
