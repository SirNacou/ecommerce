package postgres

import (
	"context"
	"fmt"

	"github.com/SirNacou/ecommerce/services/user-service/internal/app"
	"github.com/SirNacou/ecommerce/services/user-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type sqlcStore struct {
	userRepo domain.UserRepository
}

// Users implements [app.Store].
func (s *sqlcStore) Users() domain.UserRepository {
	return s.userRepo
}

type sqlcUnitOfWork struct {
	pool *pgxpool.Pool
}

// Execute implements [app.UnitOfWork].
func (uow *sqlcUnitOfWork) Execute(ctx context.Context, fn func(store app.Store) error) error {
	tx, err := uow.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	store := sqlcStore{
		userRepo: NewUserRepository(tx),
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(&store); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v, original error: %w", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
