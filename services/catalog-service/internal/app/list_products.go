package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type ListProductsQuery struct {
	PageSize   int32
	Offset     int32
	CategoryID *uuid.UUID
}

type ListProductsQueryHandler struct {
	uow UnitOfWork
}

func NewListProductsQueryHandler(uow UnitOfWork) *ListProductsQueryHandler {
	return &ListProductsQueryHandler{uow: uow}
}

func (h *ListProductsQueryHandler) Handle(ctx context.Context, q ListProductsQuery) ([]*domain.Product, error) {
	var products []*domain.Product
	err := h.uow.Execute(ctx, func(store CatalogStore) error {
		var err error
		products, err = store.ListProducts(ctx, q.CategoryID, q.PageSize, q.Offset)
		return err
	})
	return products, err
}
