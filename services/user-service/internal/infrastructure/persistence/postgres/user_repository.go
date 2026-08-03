package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/SirNacou/ecommerce/services/user-service/internal/domain"
	"github.com/SirNacou/ecommerce/services/user-service/internal/infrastructure/persistence/postgres/db"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(dbtx db.DBTX) *UserRepository {
	return &UserRepository{
		queries: db.New(dbtx),
	}
}

// Create implements [domain.UserRepository].
func (u *UserRepository) Create(ctx context.Context, user *domain.User) error {
	err := u.queries.CreateUser(ctx, db.CreateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
		CreatedAt:    user.CreatedAt,
	})

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return u.saveOutboxEvents(ctx, user)
}

func (u *UserRepository) Update(ctx context.Context, user *domain.User) error {
	err := u.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
	})

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return u.saveOutboxEvents(ctx, user)
}

// FindByEmail implements [domain.UserRepository].
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
		CreatedAt:    user.CreatedAt,
	}, nil
}

// FindByID implements [domain.UserRepository].
func (u *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := u.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
		CreatedAt:    user.CreatedAt,
	}, nil
}

// Helper method to write outbox events for both operations
func (r *UserRepository) saveOutboxEvents(ctx context.Context, user *domain.User) error {
	// events := user.PopEvents()

	// for _, event := range events {
	// 	payload, err := json.Marshal(event)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to marshal domain event: %w", err)
	// 	}

	// 	err = r.queries.CreateOutboxEvent(ctx, db.CreateOutboxEventParams{
	// 		AggregateType: "User",
	// 		AggregateID:   user.ID,
	// 		EventType:     event.EventName(),
	// 		Payload:       payload,
	// 		CreatedAt:     event.OccurredAt(),
	// 	})
	// 	if err != nil {
	// 		return fmt.Errorf("failed to save outbox event: %w", err)
	// 	}
	// }

	return nil
}
