package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SirNacou/ecommerce/services/cart-service/internal/app"
	"github.com/SirNacou/ecommerce/services/cart-service/internal/domain"
	"github.com/SirNacou/ecommerce/services/cart-service/internal/infrastructure/persistence/postgres/db"
)

type pgxUnitOfWork struct {
	pool *pgxpool.Pool
}

func NewUnitOfWork(pool *pgxpool.Pool) app.UnitOfWork {
	return &pgxUnitOfWork{pool: pool}
}

func (u *pgxUnitOfWork) Execute(ctx context.Context, fn func(store app.CartStore) error) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	store := &cartStore{
		queries: db.New(tx),
	}

	if err := fn(store); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

type cartStore struct {
	queries *db.Queries
}

func (s *cartStore) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error) {
	row, err := s.queries.GetCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCartNotFound
		}
		return nil, err
	}

	return &domain.Cart{
		ID:        row.ID,
		UserID:    row.UserID,
		Items:     make([]*domain.CartItem, 0),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (s *cartStore) CreateCart(ctx context.Context, cart *domain.Cart) error {
	return s.queries.CreateCart(ctx, db.CreateCartParams{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	})
}

func (s *cartStore) GetCartItems(ctx context.Context, cartID uuid.UUID) ([]*domain.CartItem, error) {
	rows, err := s.queries.GetCartItems(ctx, cartID)
	if err != nil {
		return nil, err
	}

	items := make([]*domain.CartItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.CartItem{
			ID:         row.ID,
			CartID:     row.CartID,
			ProductID:  row.ProductID,
			Quantity:   row.Quantity,
			PriceCents: row.PriceCents,
			CreatedAt:  row.CreatedAt,
			UpdatedAt:  row.UpdatedAt,
		})
	}
	return items, nil
}

func (s *cartStore) UpsertItem(ctx context.Context, cartID, productID uuid.UUID, quantity int32, priceCents int64) error {
	now := time.Now().UTC()
	return s.queries.UpsertCartItem(ctx, db.UpsertCartItemParams{
		ID:         uuid.New(),
		CartID:     cartID,
		ProductID:  productID,
		Quantity:   quantity,
		PriceCents: priceCents,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
}

func (s *cartStore) UpdateItemQuantity(ctx context.Context, cartID, productID uuid.UUID, quantity int32) error {
	return s.queries.UpdateCartItemQuantity(ctx, db.UpdateCartItemQuantityParams{
		CartID:    cartID,
		ProductID: productID,
		Quantity:  quantity,
		UpdatedAt: time.Now().UTC(),
	})
}

func (s *cartStore) RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error {
	return s.queries.RemoveCartItem(ctx, db.RemoveCartItemParams{
		CartID:    cartID,
		ProductID: productID,
	})
}

func (s *cartStore) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	return s.queries.ClearCartItems(ctx, cartID)
}
