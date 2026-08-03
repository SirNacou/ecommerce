package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SirNacou/ecommerce/services/catalog-service/internal/app"
	"github.com/SirNacou/ecommerce/services/catalog-service/internal/domain"
	"github.com/SirNacou/ecommerce/services/catalog-service/internal/infrastructure/persistence/postgres/db"
)

type UnitOfWork struct {
	pool *pgxpool.Pool
}

func NewUnitOfWork(pool *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{pool: pool}
}

func (u *UnitOfWork) Execute(ctx context.Context, fn func(store app.CatalogStore) error) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	queries := db.New(tx)
	store := &catalogStore{queries: queries}

	if err := fn(store); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

type catalogStore struct {
	queries *db.Queries
}

func (s *catalogStore) CreateProduct(ctx context.Context, product *domain.Product) error {
	// 1. Insert product record via SQLC
	_, err := s.queries.CreateProduct(ctx, db.CreateProductParams{
		ID:            product.ID(),
		CategoryID:    product.CategoryID(),
		Name:          product.Name(),
		Description:   product.Description(),
		PriceCents:    product.PriceCents(),
		StockQuantity: product.StockQuantity(),
	})
	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	// 2. Pop domain events and write to outbox_events table in the same transaction
	events := product.PopEvents()
	for _, event := range events {
		eventTypeGetter, ok := event.(interface{ EventType() string })
		if !ok {
			continue
		}

		payloadBytes, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal outbox event payload: %w", err)
		}

		_, err = s.queries.CreateOutboxEvent(ctx, db.CreateOutboxEventParams{
			ID:            uuid.New(),
			AggregateType: "product",
			AggregateID:   product.ID().String(),
			EventType:     eventTypeGetter.EventType(),
			Payload:       payloadBytes,
		})
		if err != nil {
			return fmt.Errorf("failed to insert outbox event: %w", err)
		}
	}

	return nil
}

func (s *catalogStore) GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	row, err := s.queries.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	return domain.ReconstituteProduct(
		row.ID,
		row.CategoryID,
		row.Name,
		row.Description,
		row.PriceCents,
		row.StockQuantity,
		row.CreatedAt,
	), nil
}

func (s *catalogStore) ListProducts(ctx context.Context, categoryID *uuid.UUID, limit int32) ([]*domain.Product, error) {
	rows, err := s.queries.ListProducts(ctx, db.ListProductsParams{
		CategoryID: categoryID,
		Limit:      limit,
	})
	if err != nil {
		return nil, err
	}

	products := make([]*domain.Product, 0, len(rows))
	for _, r := range rows {
		products = append(products, domain.ReconstituteProduct(
			r.ID,
			r.CategoryID,
			r.Name,
			r.Description,
			r.PriceCents,
			r.StockQuantity,
			r.CreatedAt,
		))
	}

	return products, nil
}
