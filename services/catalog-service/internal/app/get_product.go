package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type GetProductQueryHandler struct {
	uow UnitOfWork
}

func NewGetProductQueryHandler(uow UnitOfWork) *GetProductQueryHandler {
	return &GetProductQueryHandler{uow: uow}
}

func (h *GetProductQueryHandler) Handle(ctx context.Context, idStr string) (*domain.Product, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, domain.ErrProductNotFound
	}

	var product *domain.Product
	err = h.uow.Execute(ctx, func(store CatalogStore) error {
		var err error
		product, err = store.GetProductByID(ctx, id)
		return err
	})
	if err != nil {
		return nil, err
	}

	return product, nil
}
