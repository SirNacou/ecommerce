package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/cart-service/internal/domain"
	"github.com/google/uuid"
)

type UpdateQuantityCommand struct {
	UserID    string
	ProductID string
	Quantity  int32
}

type UpdateQuantityCommandHandler struct {
	uow        UnitOfWork
	getCartQry *GetCartQueryHandler
}

func NewUpdateQuantityCommandHandler(uow UnitOfWork, getCartQry *GetCartQueryHandler) *UpdateQuantityCommandHandler {
	return &UpdateQuantityCommandHandler{
		uow:        uow,
		getCartQry: getCartQry,
	}
}

func (h *UpdateQuantityCommandHandler) Handle(ctx context.Context, cmd UpdateQuantityCommand) (*domain.Cart, error) {
	productID, err := uuid.Parse(cmd.ProductID)
	if err != nil {
		return nil, domain.ErrItemNotFound
	}

	cart, err := h.getCartQry.Handle(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	if cmd.Quantity <= 0 {
		err = h.uow.Execute(ctx, func(store CartStore) error {
			return store.RemoveItem(ctx, cart.ID, productID)
		})
	} else {
		err = h.uow.Execute(ctx, func(store CartStore) error {
			return store.UpdateItemQuantity(ctx, cart.ID, productID, cmd.Quantity)
		})
	}

	if err != nil {
		return nil, err
	}

	return h.getCartQry.Handle(ctx, cmd.UserID)
}
