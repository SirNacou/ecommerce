package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type CreateProductInput struct {
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

func (ch *CreateProductCommandHandler) Execute(ctx context.Context, input CreateProductInput) (*domain.Product, error) {
	catID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		return nil, domain.ErrInvalidCategoryID
	}

	product, err := domain.NewProduct(catID, input.Name, input.Description, input.PriceCents, input.StockQuantity)
	if err != nil {
		return nil, err
	}

	err = ch.uow.Execute(ctx, func(store CatalogStore) error {
		return store.CreateProduct(ctx, product)
	})
	if err != nil {
		return nil, err
	}

	return product, nil
}
