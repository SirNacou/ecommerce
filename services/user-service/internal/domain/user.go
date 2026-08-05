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

type UserRegisteredEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

func (e UserRegisteredEvent) EventName() string     { return "UserRegistered" }
func (e UserRegisteredEvent) OccurredAt() time.Time { return e.Timestamp }

type User struct {
	AggregateRoot
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Name         string
	CreatedAt    time.Time
}

// NewUser enforces business validation rules during user creation.
// passwordHash is expected to already be a bcrypt hash of the raw password.
func NewUser(email, passwordHash, name string) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}
	if passwordHash == "" {
		return nil, ErrPasswordTooShort
	}

	now := time.Now().UTC()
	user := &User{
		ID:           uuid.Must(uuid.NewV7()),
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
		CreatedAt:    now,
	}

	user.AddEvent(UserRegisteredEvent{
		UserID:    user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		Timestamp: now,
	})

	return user, nil
}

func (u *User) ValidatePassword(rawPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(rawPassword))
	return err == nil
}
