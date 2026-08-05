package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SirNacou/ecommerce/services/payment-service/internal/app"
	"github.com/SirNacou/ecommerce/services/payment-service/internal/domain"
	"github.com/SirNacou/ecommerce/services/payment-service/internal/infrastructure/persistence/postgres/db"
)

type pgxUnitOfWork struct {
	pool *pgxpool.Pool
}

func NewUnitOfWork(pool *pgxpool.Pool) app.UnitOfWork {
	return &pgxUnitOfWork{pool: pool}
}

func (u *pgxUnitOfWork) Execute(ctx context.Context, fn func(store app.PaymentStore) error) error {
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

	store := &paymentStore{
		queries: db.New(tx),
	}

	if err := fn(store); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

type paymentStore struct {
	queries *db.Queries
}

func (s *paymentStore) CreatePayment(ctx context.Context, p *domain.Payment) error {
	return s.queries.CreatePayment(ctx, db.CreatePaymentParams{
		ID:            p.ID,
		OrderID:       p.OrderID,
		UserID:        p.UserID,
		AmountCents:   p.AmountCents,
		Currency:      p.Currency,
		Status:        string(p.Status),
		PaymentMethod: p.PaymentMethod,
		TransactionID: p.TransactionID,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	})
}

func (s *paymentStore) GetPaymentByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	row, err := s.queries.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, err
	}

	return &domain.Payment{
		ID:            row.ID,
		OrderID:       row.OrderID,
		UserID:        row.UserID,
		AmountCents:   row.AmountCents,
		Currency:      row.Currency,
		Status:        domain.PaymentStatus(row.Status),
		PaymentMethod: row.PaymentMethod,
		TransactionID: row.TransactionID,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}, nil
}

func (s *paymentStore) GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error) {
	row, err := s.queries.GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, err
	}

	return &domain.Payment{
		ID:            row.ID,
		OrderID:       row.OrderID,
		UserID:        row.UserID,
		AmountCents:   row.AmountCents,
		Currency:      row.Currency,
		Status:        domain.PaymentStatus(row.Status),
		PaymentMethod: row.PaymentMethod,
		TransactionID: row.TransactionID,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}, nil
}

func (s *paymentStore) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus) error {
	return s.queries.UpdatePaymentStatus(ctx, db.UpdatePaymentStatusParams{
		ID:        id,
		Status:    string(status),
		UpdatedAt: time.Now().UTC(),
	})
}

func (s *paymentStore) SaveOutboxEvent(ctx context.Context, aggregateType, aggregateID, eventType string, payload []byte) error {
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
