package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/cart-service/internal/domain"
	"github.com/google/uuid"
)

type RemoveItemCommandHandler struct {
	uow        UnitOfWork
	getCartQry *GetCartQueryHandler
}

func NewRemoveItemCommandHandler(uow UnitOfWork, getCartQry *GetCartQueryHandler) *RemoveItemCommandHandler {
	return &RemoveItemCommandHandler{
		uow:        uow,
		getCartQry: getCartQry,
	}
}

func (h *RemoveItemCommandHandler) Handle(ctx context.Context, userIDStr, productIDStr string) (*domain.Cart, error) {
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return nil, domain.ErrItemNotFound
	}

	cart, err := h.getCartQry.Handle(ctx, userIDStr)
	if err != nil {
		return nil, err
	}

	err = h.uow.Execute(ctx, func(store CartStore) error {
		return store.RemoveItem(ctx, cart.ID, productID)
	})
	if err != nil {
		return nil, err
	}

	return h.getCartQry.Handle(ctx, userIDStr)
}
