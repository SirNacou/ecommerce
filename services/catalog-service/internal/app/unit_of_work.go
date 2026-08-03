package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type CatalogStore interface {
	CreateProduct(ctx context.Context, product *domain.Product) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	ListProducts(ctx context.Context, categoryID *uuid.UUID, limit int32) ([]*domain.Product, error)
}

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(store CatalogStore) error) error
}