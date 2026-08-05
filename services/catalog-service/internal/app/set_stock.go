package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type SetStockCommand struct {
	ProductID string
	Quantity  int32
}

type SetStockCommandHandler struct {
	uow UnitOfWork
}

func NewSetStockCommandHandler(uow UnitOfWork) *SetStockCommandHandler {
	return &SetStockCommandHandler{uow: uow}
}

func (h *SetStockCommandHandler) Handle(ctx context.Context, cmd SetStockCommand) (*domain.InventoryItem, error) {
	productID, err := uuid.Parse(cmd.ProductID)
	if err != nil {
		return nil, domain.ErrItemNotFound
	}

	item, err := domain.NewInventoryItem(productID, cmd.Quantity)
	if err != nil {
		return nil, err
	}

	err = h.uow.Execute(ctx, func(store CatalogStore) error {
		return store.UpsertItem(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	return item, nil
}
