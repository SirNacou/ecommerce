package app

import (
	"context"
	"errors"

	"github.com/SirNacou/ecommerce/services/user-service/internal/domain"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type LoginUserCommand struct {
	Email    string
	Password string
}

type LoginUserResult struct {
	AccessToken  string
	RefreshToken string
}

type LoginUserUseCase struct {
	uow           UnitOfWork
	hasher        PasswordHasher
	tokenProvider TokenProvider
}

func NewLoginUserUseCase(
	uow UnitOfWork,
	hasher PasswordHasher,
	tokenProvider TokenProvider,
) *LoginUserUseCase {
	return &LoginUserUseCase{
		uow:           uow,
		hasher:        hasher,
		tokenProvider: tokenProvider,
	}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, cmd LoginUserCommand) (*LoginUserResult, error) {
	var user *domain.User

	err := uc.uow.Execute(ctx, func(store Store) error {
		u, err := store.Users().FindByEmail(ctx, cmd.Email)
		if err != nil {
			return ErrInvalidCredentials
		}
		user = u
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Compare hashed password
	if !uc.hasher.Compare(user.PasswordHash, cmd.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT tokens
	accToken, refToken, err := uc.tokenProvider.GenerateTokens(user.ID.String())
	if err != nil {
		return nil, err
	}

	return &LoginUserResult{
		AccessToken:  accToken,
		RefreshToken: refToken,
	}, nil
}
