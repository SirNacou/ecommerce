package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type GetProductsByIdsQueryHandler struct {
	uow UnitOfWork
}

func NewGetProductsByIdsQueryHandler(uow UnitOfWork) *GetProductsByIdsQueryHandler {
	return &GetProductsByIdsQueryHandler{uow: uow}
}

func (h *GetProductsByIdsQueryHandler) Handle(ctx context.Context, idStrs []string) ([]*domain.Product, error) {
	ids := make([]uuid.UUID, 0, len(idStrs))
	for _, s := range idStrs {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, domain.ErrProductNotFound
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return []*domain.Product{}, nil
	}

	var products []*domain.Product
	err := h.uow.Execute(ctx, func(store CatalogStore) error {
		var err error
		products, err = store.ListProductsByIds(ctx, ids)
		return err
	})
	if err != nil {
		return nil, err
	}

	return products, nil
}