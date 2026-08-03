package app

// PasswordHasher defines the port for hashing and verifying credentials
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

// TokenProvider defines the port for generating and validating auth tokens
type TokenProvider interface {
	GenerateTokens(userID, email string) (accessToken string, refreshToken string, err error)
	ValidateToken(tokenStr string) (claims *UserClaims, err error)
}

type UserClaims struct {
	UserID string
	Email  string
}
