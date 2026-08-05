package postgres

import (
	"context"

	"github.com/SirNacou/ecommerce/pkg/eventbus"
	"github.com/SirNacou/ecommerce/services/user-service/internal/infrastructure/persistence/postgres/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OutboxReader adapts the sqlc outbox queries to the eventbus.OutboxReader
// interface for the outbox dispatcher.
type OutboxReader struct {
	queries *db.Queries
}

func NewOutboxReader(pool *pgxpool.Pool) *OutboxReader {
	return &OutboxReader{queries: db.New(pool)}
}

func (r *OutboxReader) GetPendingOutboxEvents(ctx context.Context, limit int32) ([]eventbus.OutboxEvent, error) {
	rows, err := r.queries.GetPendingOutboxEvents(ctx, limit)
	if err != nil {
		return nil, err
	}

	events := make([]eventbus.OutboxEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, eventbus.OutboxEvent{
			ID:        row.ID.String(),
			EventType: row.EventType,
			Payload:   row.Payload,
		})
	}

	return events, nil
}

func (r *OutboxReader) MarkOutboxEventProcessed(ctx context.Context, id string) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.queries.MarkOutboxEventProcessed(ctx, parsed)
}