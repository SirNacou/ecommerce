package grpc

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"github.com/SirNacou/ecommerce/services/user-service/internal/app"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

// NewAuthInterceptor validates JWT tokens for protected RPC endpoints
func NewAuthInterceptor(tokenProvider app.TokenProvider) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			procedure := req.Spec().Procedure

			// 1. Whitelist public endpoints (Register & Login do not require tokens)
			if procedure == "/user.v1.UserService/Register" || procedure == "/user.v1.UserService/Login" {
				return next(ctx, req)
			}

			// 2. Extract Authorization header
			authHeader := req.Header().Get("Authorization")
			if authHeader == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing authorization header"))
			}

			// 3. Parse "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid authorization header format"))
			}

			tokenString := parts[1]

			// 4. Validate JWT Token using your JWT TokenProvider
			claims, err := tokenProvider.ValidateToken(tokenString)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid or expired token"))
			}

			// 5. Inject UserID into request Context
			ctx = context.WithValue(ctx, UserIDContextKey, claims.UserID)

			return next(ctx, req)
		}
	}
}

// GetUserIDFromContext retrieves the authenticated User ID inside RPC handlers
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	return userID, ok
}