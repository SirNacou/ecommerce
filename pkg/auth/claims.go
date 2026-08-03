package auth

import "context"

type contextKey string

const claimsContextKey contextKey = "user_claims"

// UserClaims holds authenticated user identity & RBAC claims
type UserClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// WithClaims injects UserClaims into context.Context
func WithClaims(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

// GetClaims retrieves UserClaims from context.Context
func GetClaims(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*UserClaims)
	return claims, ok
}

// GetUserID retrieves the authenticated User ID directly from context.Context
func GetUserID(ctx context.Context) (string, bool) {
	if claims, ok := GetClaims(ctx); ok {
		return claims.UserID, true
	}
	return "", false
}
