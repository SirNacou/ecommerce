package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	// Create persists a new user in the repository
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	// FindByEmail retrieves a user by their email address
	FindByEmail(ctx context.Context, email string) (*User, error)
	// FindByID retrieves a user by their ID
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
}
