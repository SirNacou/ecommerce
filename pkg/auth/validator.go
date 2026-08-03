package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type TokenValidator interface {
	VerifyToken(tokenString string) (*UserClaims, error)
}

type JWTValidator struct {
	secretKey []byte
}

func NewJWTValidator(secretKey string) *JWTValidator {
	return &JWTValidator{secretKey: []byte(secretKey)}
}

type customClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email,omitempty"`
	TokenType string `json:"token_type"` // Must be "access"
	jwt.RegisteredClaims
}

func (v *JWTValidator) VerifyToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return v.secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok {
		return nil, errors.New("invalid token claims format")
	}

	// Reject Refresh Tokens sent to API endpoints
	if claims.TokenType != "access" {
		return nil, errors.New("invalid token type for authorization")
	}

	return &UserClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
	}, nil
}
