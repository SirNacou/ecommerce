package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/catalog-service/internal/domain"
)

type CreateCategoryCommand struct {
	Name string
	Slug string
}

type CreateCategoryCommandHandler struct {
	uow UnitOfWork
}

func NewCreateCategoryCommandHandler(uow UnitOfWork) *CreateCategoryCommandHandler {
	return &CreateCategoryCommandHandler{uow: uow}
}

func (h *CreateCategoryCommandHandler) Handle(ctx context.Context, cmd CreateCategoryCommand) (*domain.Category, error) {
	category, err := domain.NewCategory(cmd.Name, cmd.Slug)
	if err != nil {
		return nil, err
	}

	err = h.uow.Execute(ctx, func(store CatalogStore) error {
		return store.CreateCategory(ctx, category)
	})
	if err != nil {
		return nil, err
	}

	return category, nil
}
