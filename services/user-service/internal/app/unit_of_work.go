package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/user-service/internal/domain"
)

type Store interface {
	Users() domain.UserRepository
}

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(store Store) error) error
}
