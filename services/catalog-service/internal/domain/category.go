package domain

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

func NewCategory(name, slug string) (*Category, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	return &Category{
		ID:        uuid.New(),
		Name:      name,
		Slug:      slug,
		CreatedAt: time.Now().UTC(),
	}, nil
}
