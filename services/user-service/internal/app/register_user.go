package app

import (
	"context"
	"fmt"

	"github.com/SirNacou/ecommerce/services/user-service/internal/domain"
)

type RegisterUserCommand struct {
	Email    string
	Password string
	Name     string
}

type RegisterUserResult struct {
	User *struct {
		ID    string
		Email string
		Name  string
	}
}

type RegisterUserUseCase struct {
	uow    UnitOfWork
	hasher PasswordHasher
}

func NewRegisterUserUseCase(uow UnitOfWork, hasher PasswordHasher) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		uow:    uow,
		hasher: hasher,
	}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, cmd RegisterUserCommand) (*RegisterUserResult, error) {
	if len(cmd.Password) < 6 {
		return nil, domain.ErrPasswordTooShort
	}

	hashedPassword, err := uc.hasher.Hash(cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	var createdUser *domain.User

	// Execute inside transactional Unit of Work
	err = uc.uow.Execute(ctx, func(store Store) error {
		// NewUser creates aggregate root & records UserRegisteredEvent
		user, err := domain.NewUser(cmd.Email, hashedPassword, cmd.Name)
		if err != nil {
			return err
		}

		// Explicit Create call saves user + flushes outbox events
		if err := store.Users().Create(ctx, user); err != nil {
			return err
		}

		createdUser = user
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &RegisterUserResult{
		User: &struct {
			ID    string
			Email string
			Name  string
		}{
			ID:    createdUser.ID.String(),
			Email: createdUser.Email,
			Name:  createdUser.Name,
		},
	}, nil
}
