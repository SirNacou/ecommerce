package security

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	return &BcryptHasher{
		cost: cost,
	}
}

// Compare implements [app.PasswordHasher].
func (b *BcryptHasher) Compare(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Hash implements [app.PasswordHasher].
func (b *BcryptHasher) Hash(password string) (string, error) {
	bytePassword, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}

	return string(bytePassword), nil
}
