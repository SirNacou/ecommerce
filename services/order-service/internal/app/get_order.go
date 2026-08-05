package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/order-service/internal/domain"
	"github.com/google/uuid"
)

type GetOrderQueryHandler struct {
	uow UnitOfWork
}

func NewGetOrderQueryHandler(uow UnitOfWork) *GetOrderQueryHandler {
	return &GetOrderQueryHandler{uow: uow}
}

func (h *GetOrderQueryHandler) Handle(ctx context.Context, orderIDStr, userIDStr string) (*domain.Order, error) {
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}

	var order *domain.Order
	err = h.uow.Execute(ctx, func(store OrderStore) error {
		var err error
		order, err = store.GetOrderByID(ctx, orderID)
		return err
	})
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, domain.ErrOrderNotFound
	}

	return order, nil
}
