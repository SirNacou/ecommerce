package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/user-service/internal/domain"
	"github.com/google/uuid"
)

type GetUserQuery struct {
	ID uuid.UUID
}

type GetUserResult struct {
	User *domain.User
}

type GetUserQueryHandler struct {
	uow UnitOfWork
}

func NewGetUserQueryHandler(uow UnitOfWork) *GetUserQueryHandler {
	return &GetUserQueryHandler{uow: uow}
}

func (uc *GetUserQueryHandler) Execute(ctx context.Context, query GetUserQuery) (*GetUserResult, error) {
	var user *domain.User

	err := uc.uow.Execute(ctx, func(store Store) error {
		u, err := store.Users().FindByID(ctx, query.ID)
		if err != nil {
			return err
		}
		user = u
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &GetUserResult{User: user}, nil
}
