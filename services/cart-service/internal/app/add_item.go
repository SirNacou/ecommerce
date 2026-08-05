package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/cart-service/internal/domain"
	"github.com/google/uuid"
)

type AddItemCommand struct {
	UserID    string
	ProductID string
	Quantity  int32
}

// PriceResolver resolves authoritative product prices for a set of product ids.
type PriceResolver interface {
	ResolvePrices(ctx context.Context, productIDs []string) (map[string]int64, error)
}

type AddItemCommandHandler struct {
	uow        UnitOfWork
	getCartQry *GetCartQueryHandler
	prices     PriceResolver
}

func NewAddItemCommandHandler(uow UnitOfWork, getCartQry *GetCartQueryHandler, prices PriceResolver) *AddItemCommandHandler {
	return &AddItemCommandHandler{
		uow:        uow,
		getCartQry: getCartQry,
		prices:     prices,
	}
}

func (h *AddItemCommandHandler) Handle(ctx context.Context, cmd AddItemCommand) (*domain.Cart, error) {
	if cmd.Quantity <= 0 {
		return nil, domain.ErrInvalidQuantity
	}

	productID, err := uuid.Parse(cmd.ProductID)
	if err != nil {
		return nil, domain.ErrItemNotFound
	}

	prices, err := h.prices.ResolvePrices(ctx, []string{cmd.ProductID})
	if err != nil {
		return nil, err
	}
	priceCents, ok := prices[cmd.ProductID]
	if !ok {
		return nil, domain.ErrProductNotFound
	}

	cart, err := h.getCartQry.Handle(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	err = h.uow.Execute(ctx, func(store CartStore) error {
		return store.UpsertItem(ctx, cart.ID, productID, cmd.Quantity, priceCents)
	})
	if err != nil {
		return nil, err
	}

	return h.getCartQry.Handle(ctx, cmd.UserID)
}
