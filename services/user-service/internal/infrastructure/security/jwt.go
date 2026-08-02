package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type Claims struct {
	UserID string `json:"user_id"`
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
func (j *JWTProvider) GenerateTokens(userID string) (accessToken string, refreshToken string, err error) {
	now := time.Now().UTC()

	accessClaims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   "",
			Audience:  jwt.ClaimStrings{},
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessExpiry)),
			NotBefore: &jwt.NumericDate{},
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        userID,
		},
	}

	accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accToken.SignedString(j.secretKey)
	if err != nil {
		return "", "", err
	}

	refreshClaims := Claims{
		UserID: userID,
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
func (j *JWTProvider) ValidateToken(tokenStr string) (userID string, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secretKey, nil
	})

	if err != nil || !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", ErrInvalidToken
	}

	return claims.UserID, nil
}
