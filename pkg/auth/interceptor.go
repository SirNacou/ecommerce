package auth

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
)

// NewConnectInterceptor creates a reusable ConnectRPC authentication middleware.
// Pass publicProcedures map to mark specific RPC methods as unauthenticated/public.
func NewConnectInterceptor(validator TokenValidator, publicProcedures map[string]bool) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			procedure := req.Spec().Procedure

			// 1. Skip token check if procedure is marked as public
			if publicProcedures != nil && publicProcedures[procedure] {
				return next(ctx, req)
			}

			// 2. Extract Authorization header
			authHeader := req.Header().Get("Authorization")
			if authHeader == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing authorization header"))
			}

			// 3. Parse "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid authorization header format"))
			}

			// 4. Validate Token
			claims, err := validator.VerifyToken(parts[1])
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			// 5. Inject UserClaims into request context
			ctx = WithClaims(ctx, claims)

			return next(ctx, req)
		}
	}
}
