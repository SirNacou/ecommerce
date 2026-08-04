package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/catalog-service/internal/domain"
	"github.com/google/uuid"
)

type CatalogStore interface {
	CreateProduct(ctx context.Context, product *domain.Product) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	ListProducts(ctx context.Context, categoryID *uuid.UUID, limit, offset int32) ([]*domain.Product, error)
	CreateCategory(ctx context.Context, category *domain.Category) error
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	ListCategories(ctx context.Context) ([]*domain.Category, error)
}

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(store CatalogStore) error) error
}
