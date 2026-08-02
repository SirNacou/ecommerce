package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail     = errors.New("invalid email address")
	ErrPasswordTooShort = errors.New("password must be at least 6 characters")
)

type User struct {
	AggregateRoot
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Name         string
	CreatedAt    time.Time
}

// NewUser enforces business validation rules during user creation
func NewUser(email, rawPassword, name string) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}
	if len(rawPassword) < 6 {
		return nil, ErrPasswordTooShort
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           uuid.Must(uuid.NewV7()),
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
		CreatedAt:    time.Now().UTC(),
	}, nil
}

func (u *User) ValidatePassword(rawPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(rawPassword))
	return err == nil
}
