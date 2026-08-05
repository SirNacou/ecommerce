package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type GetStockQueryHandler struct {
	uow UnitOfWork
}

func NewGetStockQueryHandler(uow UnitOfWork) *GetStockQueryHandler {
	return &GetStockQueryHandler{uow: uow}
}

func (h *GetStockQueryHandler) Handle(ctx context.Context, productIDStr string) (*domain.InventoryItem, error) {
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return nil, domain.ErrItemNotFound
	}

	var item *domain.InventoryItem
	err = h.uow.Execute(ctx, func(store CatalogStore) error {
		var err error
		item, err = store.GetItem(ctx, productID)
		return err
	})

	return item, err
}
