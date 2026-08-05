package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/cart-service/internal/domain"
	"github.com/google/uuid"
)

type GetCartQueryHandler struct {
	uow UnitOfWork
}

func NewGetCartQueryHandler(uow UnitOfWork) *GetCartQueryHandler {
	return &GetCartQueryHandler{uow: uow}
}

func (h *GetCartQueryHandler) Handle(ctx context.Context, userIDStr string) (*domain.Cart, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, domain.ErrCartNotFound
	}

	var cart *domain.Cart
	err = h.uow.Execute(ctx, func(store CartStore) error {
		var err error
		cart, err = store.GetCartByUserID(ctx, userID)
		if err != nil && err == domain.ErrCartNotFound {
			cart = domain.NewCart(userID)
			if createErr := store.CreateCart(ctx, cart); createErr != nil {
				return createErr
			}
		} else if err != nil {
			return err
		}

		items, err := store.GetCartItems(ctx, cart.ID)
		if err != nil {
			return err
		}
		cart.Items = items
		return nil
	})

	return cart, err
}
