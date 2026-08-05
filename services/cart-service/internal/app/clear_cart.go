package app

import (
	"context"
)

type ClearCartCommandHandler struct {
	uow        UnitOfWork
	getCartQry *GetCartQueryHandler
}

func NewClearCartCommandHandler(uow UnitOfWork, getCartQry *GetCartQueryHandler) *ClearCartCommandHandler {
	return &ClearCartCommandHandler{
		uow:        uow,
		getCartQry: getCartQry,
	}
}

func (h *ClearCartCommandHandler) Handle(ctx context.Context, userIDStr string) error {
	cart, err := h.getCartQry.Handle(ctx, userIDStr)
	if err != nil {
		return err
	}

	return h.uow.Execute(ctx, func(store CartStore) error {
		return store.ClearCart(ctx, cart.ID)
	})
}
