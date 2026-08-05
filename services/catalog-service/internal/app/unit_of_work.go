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
	ListProductsByIds(ctx context.Context, ids []uuid.UUID) ([]*domain.Product, error)
	CreateCategory(ctx context.Context, category *domain.Category) error
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	ListCategories(ctx context.Context) ([]*domain.Category, error)
	GetItemForUpdate(ctx context.Context, productID uuid.UUID) (*domain.InventoryItem, error)
	GetItem(ctx context.Context, productID uuid.UUID) (*domain.InventoryItem, error)
	UpsertItem(ctx context.Context, item *domain.InventoryItem) error
	UpdateStock(ctx context.Context, item *domain.InventoryItem) error
	CreateReservation(ctx context.Context, reservationID, productID uuid.UUID, quantity int32) error
	GetReservation(ctx context.Context, reservationID uuid.UUID) (*domain.StockReservation, error)
	UpdateReservationStatus(ctx context.Context, reservationID uuid.UUID, status string) error
	SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload []byte) error
}

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(store CatalogStore) error) error
}
