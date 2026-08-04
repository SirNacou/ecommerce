package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/catalog-service/internal/domain"
)

type ListCategoriesQueryHandler struct {
	uow UnitOfWork
}

func NewListCategoriesQueryHandler(uow UnitOfWork) *ListCategoriesQueryHandler {
	return &ListCategoriesQueryHandler{uow: uow}
}

func (h *ListCategoriesQueryHandler) Handle(ctx context.Context) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := h.uow.Execute(ctx, func(store CatalogStore) error {
		var err error
		categories, err = store.ListCategories(ctx)
		return err
	})
	return categories, err
}
