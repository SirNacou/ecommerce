package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SirNacou/ecommerce/services/notification-service/internal/app"
	"github.com/SirNacou/ecommerce/services/notification-service/internal/domain"
	"github.com/SirNacou/ecommerce/services/notification-service/internal/infrastructure/persistence/postgres/db"
)

type pgxUnitOfWork struct {
	pool *pgxpool.Pool
}

func NewUnitOfWork(pool *pgxpool.Pool) app.UnitOfWork {
	return &pgxUnitOfWork{pool: pool}
}

func (u *pgxUnitOfWork) Execute(ctx context.Context, fn func(store app.NotificationStore) error) error {
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

	store := &notificationStore{
		queries: db.New(tx),
	}

	if err := fn(store); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

type notificationStore struct {
	queries *db.Queries
}

func (s *notificationStore) CreateNotification(ctx context.Context, n *domain.Notification) error {
	return s.queries.CreateNotification(ctx, db.CreateNotificationParams{
		ID:        n.ID,
		UserID:    n.UserID,
		Channel:   string(n.Channel),
		Recipient: n.Recipient,
		Subject:   n.Subject,
		Body:      n.Body,
		Status:    string(n.Status),
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	})
}

func (s *notificationStore) GetNotificationByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	row, err := s.queries.GetNotificationByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotificationNotFound
		}
		return nil, err
	}

	return &domain.Notification{
		ID:        row.ID,
		UserID:    row.UserID,
		Channel:   domain.Channel(row.Channel),
		Recipient: row.Recipient,
		Subject:   row.Subject,
		Body:      row.Body,
		Status:    domain.NotificationStatus(row.Status),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (s *notificationStore) UpdateNotificationStatus(ctx context.Context, id uuid.UUID, status domain.NotificationStatus) error {
	return s.queries.UpdateNotificationStatus(ctx, db.UpdateNotificationStatusParams{
		ID:        id,
		Status:    string(status),
		UpdatedAt: time.Now().UTC(),
	})
}

func (s *notificationStore) SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload []byte) error {
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
