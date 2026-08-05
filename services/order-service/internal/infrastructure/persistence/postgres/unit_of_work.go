package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SirNacou/ecommerce/services/order-service/internal/app"
	"github.com/SirNacou/ecommerce/services/order-service/internal/domain"
	"github.com/SirNacou/ecommerce/services/order-service/internal/infrastructure/persistence/postgres/db"
)

type pgxUnitOfWork struct {
	pool *pgxpool.Pool
}

func NewUnitOfWork(pool *pgxpool.Pool) app.UnitOfWork {
	return &pgxUnitOfWork{pool: pool}
}

func (u *pgxUnitOfWork) Execute(ctx context.Context, fn func(store app.OrderStore) error) error {
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

	store := &orderStore{
		queries: db.New(tx),
	}

	if err := fn(store); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

type orderStore struct {
	queries *db.Queries
}

func (s *orderStore) CreateOrder(ctx context.Context, order *domain.Order) error {
	err := s.queries.CreateOrder(ctx, db.CreateOrderParams{
		ID:         order.ID,
		UserID:     order.UserID,
		TotalCents: order.TotalCents,
		Status:     string(order.Status),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	})
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		err := s.queries.CreateOrderItem(ctx, db.CreateOrderItemParams{
			ID:         item.ID,
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			PriceCents: item.PriceCents,
			CreatedAt:  item.CreatedAt,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *orderStore) GetOrderByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	row, err := s.queries.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	itemRows, err := s.queries.GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return nil, err
	}

	items := make([]*domain.OrderItem, 0, len(itemRows))
	for _, itemRow := range itemRows {
		items = append(items, &domain.OrderItem{
			ID:         itemRow.ID,
			OrderID:    itemRow.OrderID,
			ProductID:  itemRow.ProductID,
			Quantity:   itemRow.Quantity,
			PriceCents: itemRow.PriceCents,
			CreatedAt:  itemRow.CreatedAt,
		})
	}

	return &domain.Order{
		ID:         row.ID,
		UserID:     row.UserID,
		Items:      items,
		TotalCents: row.TotalCents,
		Status:     domain.OrderStatus(row.Status),
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}, nil
}

func (s *orderStore) ListOrdersByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*domain.Order, error) {
	rows, err := s.queries.ListOrdersByUserID(ctx, db.ListOrdersByUserIDParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	orders := make([]*domain.Order, 0, len(rows))
	for _, row := range rows {
		itemRows, err := s.queries.GetOrderItemsByOrderID(ctx, row.ID)
		if err != nil {
			return nil, err
		}

		items := make([]*domain.OrderItem, 0, len(itemRows))
		for _, itemRow := range itemRows {
			items = append(items, &domain.OrderItem{
				ID:         itemRow.ID,
				OrderID:    itemRow.OrderID,
				ProductID:  itemRow.ProductID,
				Quantity:   itemRow.Quantity,
				PriceCents: itemRow.PriceCents,
				CreatedAt:  itemRow.CreatedAt,
			})
		}

		orders = append(orders, &domain.Order{
			ID:         row.ID,
			UserID:     row.UserID,
			Items:      items,
			TotalCents: row.TotalCents,
			Status:     domain.OrderStatus(row.Status),
			CreatedAt:  row.CreatedAt,
			UpdatedAt:  row.UpdatedAt,
		})
	}

	return orders, nil
}

func (s *orderStore) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	return s.queries.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		ID:        id,
		Status:    string(status),
		UpdatedAt: time.Now().UTC(),
	})
}

func (s *orderStore) SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload []byte) error {
	return s.queries.InsertOutboxEvent(ctx, db.InsertOutboxEventParams{
		ID:            uuid.New(),
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		Payload:       payload,
		Status:        "PENDING",
		CreatedAt:     time.Now().UTC(),
	})
}
