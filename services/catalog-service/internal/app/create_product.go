package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type CreateProductCommand struct {
	CategoryID    string
	Name          string
	Description   string
	PriceCents    int64
	StockQuantity int32
}

type CreateProductCommandHandler struct {
	uow UnitOfWork
}

func NewCreateProductCommandHandler(uow UnitOfWork) *CreateProductCommandHandler {
	return &CreateProductCommandHandler{uow: uow}
}

func (h *CreateProductCommandHandler) Handle(ctx context.Context, cmd CreateProductCommand) (*domain.Product, error) {
	categoryID, err := uuid.Parse(cmd.CategoryID)
	if err != nil {
		return nil, domain.ErrCategoryNotFound
	}

	product, err := domain.NewProduct(
		categoryID,
		cmd.Name,
		cmd.Description,
		cmd.PriceCents,
		cmd.StockQuantity,
	)
	if err != nil {
		return nil, err
	}

	err = h.uow.Execute(ctx, func(store CatalogStore) error {
		return store.CreateProduct(ctx, product)
	})
	if err != nil {
		return nil, err
	}

	return product, nil
}
