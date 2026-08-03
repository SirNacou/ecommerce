package grpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	userv1 "github.com/SirNacou/ecommerce/services/user-service/gen/v1"
	"github.com/SirNacou/ecommerce/services/user-service/gen/v1/userv1connect"
	"github.com/SirNacou/ecommerce/services/user-service/internal/app"
	"github.com/SirNacou/ecommerce/services/user-service/internal/domain"
	"github.com/SirNacou/ecommerce/services/user-service/internal/infrastructure/persistence/postgres"
)

type UserHandler struct {
	userv1connect.UnimplementedUserServiceHandler
	registerHandler *app.RegisterUserCommandHandler
	loginHandler    *app.LoginUserCommandHandler
	getUserHandler  *app.GetUserQueryHandler
}

func NewUserHandler(
	registerHandler *app.RegisterUserCommandHandler,
	loginHandler *app.LoginUserCommandHandler,
	getUserHandler *app.GetUserQueryHandler,
) *UserHandler {
	return &UserHandler{
		registerHandler: registerHandler,
		loginHandler:    loginHandler,
		getUserHandler:  getUserHandler,
	}
}

func (h *UserHandler) Register(
	ctx context.Context,
	req *connect.Request[userv1.RegisterRequest],
) (*connect.Response[userv1.RegisterResponse], error) {
	res, err := h.registerHandler.Execute(ctx, app.RegisterUserCommand{
		Email:    req.Msg.GetEmail(),
		Password: req.Msg.GetPassword(),
		Name:     req.Msg.GetName(),
	})
	if err != nil {
		return nil, mapToConnectError(err)
	}

	return connect.NewResponse(&userv1.RegisterResponse{
		Id:    res.User.ID,
		Email: res.User.Email,
	}), nil
}

func (h *UserHandler) Login(
	ctx context.Context,
	req *connect.Request[userv1.LoginRequest],
) (*connect.Response[userv1.LoginResponse], error) {
	res, err := h.loginHandler.Execute(ctx, app.LoginUserCommand{
		Email:    req.Msg.GetEmail(),
		Password: req.Msg.GetPassword(),
	})
	if err != nil {
		return nil, mapToConnectError(err)
	}

	return connect.NewResponse(&userv1.LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}), nil
}

// Map domain and application errors to ConnectRPC Status Codes
func mapToConnectError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidEmail), errors.Is(err, domain.ErrPasswordTooShort):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, app.ErrInvalidCredentials):
		return connect.NewError(connect.CodeUnauthenticated, err)
	case errors.Is(err, postgres.ErrUserNotFound):
		return connect.NewError(connect.CodeNotFound, err)
	default:
		return connect.NewError(connect.CodeInternal, errors.New("internal server error"))
	}
}
