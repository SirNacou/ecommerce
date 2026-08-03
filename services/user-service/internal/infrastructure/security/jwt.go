package security

import (
	"errors"
	"time"

	"github.com/SirNacou/ecommerce/services/user-service/internal/app"
	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type CustomClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email,omitempty"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}
type JWTProvider struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTProvider(secretKey []byte, accessExpiry, refreshExpiry time.Duration) *JWTProvider {
	return &JWTProvider{
		secretKey:     secretKey,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateTokens implements [app.TokenProvider].
func (j *JWTProvider) GenerateTokens(userID, email string) (accessToken string, refreshToken string, err error) {
	now := time.Now().UTC()

	accessClaims := CustomClaims{
		UserID:    userID,
		Email:     email,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        userID,
		},
	}

	accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accToken.SignedString(j.secretKey)
	if err != nil {
		return "", "", err
	}

	refreshClaims := CustomClaims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        userID,
		},
	}
	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refToken.SignedString(j.secretKey)
	if err != nil {
		return "", "", err
	}

	return accessStr, refreshStr, nil
}

// ValidateToken implements [app.TokenProvider].
func (j *JWTProvider) ValidateToken(tokenStr string) (*app.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return &app.UserClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
	}, nil
}
