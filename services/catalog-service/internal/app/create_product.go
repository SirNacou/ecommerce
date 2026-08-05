package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type CreateProductCommand struct {
	CategoryID  string
	Name        string
	Description string
	PriceCents  int64
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
	)
	if err != nil {
		return nil, err
	}

	item, err := domain.NewInventoryItem(product.ID, 0)
	if err != nil {
		return nil, err
	}

	err = h.uow.Execute(ctx, func(store CatalogStore) error {
		if err := store.CreateProduct(ctx, product); err != nil {
			return err
		}
		return store.UpsertItem(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	return product, nil
}
