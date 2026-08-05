package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/order-service/internal/domain"
	"github.com/google/uuid"
)

type ListOrdersQuery struct {
	UserID   string
	PageSize int32
	Offset   int32
}

type ListOrdersQueryHandler struct {
	uow UnitOfWork
}

func NewListOrdersQueryHandler(uow UnitOfWork) *ListOrdersQueryHandler {
	return &ListOrdersQueryHandler{uow: uow}
}

func (h *ListOrdersQueryHandler) Handle(ctx context.Context, q ListOrdersQuery) ([]*domain.Order, error) {
	userID, err := uuid.Parse(q.UserID)
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}

	var orders []*domain.Order
	err = h.uow.Execute(ctx, func(store OrderStore) error {
		var err error
		orders, err = store.ListOrdersByUserID(ctx, userID, q.PageSize, q.Offset)
		return err
	})

	return orders, err
}
