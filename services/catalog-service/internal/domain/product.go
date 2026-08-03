package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPrice = errors.New("price must be greater than or equal to zero")
	ErrInvalidName  = errors.New("product name cannot be empty")
	ErrProductNotFound = errors.New("product not found")
)

type ProductCreatedEvent struct {
	ProductID     string    `json:"product_id"`
	CategoryID    string    `json:"category_id"`
	Name          string    `json:"name"`
	PriceCents    int64     `json:"price_cents"`
	StockQuantity int32     `json:"stock_quantity"`
	OccurredAt    time.Time `json:"occurred_at"`
}

func (e ProductCreatedEvent) EventType() string { return "catalog.product.created" }

type Product struct {
	id            uuid.UUID
	categoryID    uuid.UUID
	name          string
	description   string
	priceCents    int64
	stockQuantity int32
	createdAt     time.Time
	pendingEvents []any
}

func NewProduct(categoryID uuid.UUID, name, description string, priceCents int64, stock int32) (*Product, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	if priceCents < 0 {
		return nil, ErrInvalidPrice
	}

	pID := uuid.New()
	p := &Product{
		id:            pID,
		categoryID:    categoryID,
		name:          name,
		description:   description,
		priceCents:    priceCents,
		stockQuantity: stock,
		createdAt:     time.Now().UTC(),
	}

	// Record Domain Event for Outbox Pattern
	p.pendingEvents = append(p.pendingEvents, ProductCreatedEvent{
		ProductID:     pID.String(),
		CategoryID:    categoryID.String(),
		Name:          name,
		PriceCents:    priceCents,
		StockQuantity: stock,
		OccurredAt:    p.createdAt,
	})

	return p, nil
}

func ReconstituteProduct(id, categoryID uuid.UUID, name, description string, priceCents int64, stock int32, createdAt time.Time) *Product {
	return &Product{
		id:            id,
		categoryID:    categoryID,
		name:          name,
		description:   description,
		priceCents:    priceCents,
		stockQuantity: stock,
		createdAt:     createdAt,
	}
}

func (p *Product) ID() uuid.UUID         { return p.id }
func (p *Product) CategoryID() uuid.UUID { return p.categoryID }
func (p *Product) Name() string          { return p.name }
func (p *Product) Description() string   { return p.description }
func (p *Product) PriceCents() int64     { return p.priceCents }
func (p *Product) StockQuantity() int32  { return p.stockQuantity }
func (p *Product) CreatedAt() time.Time  { return p.createdAt }

func (p *Product) PopEvents() []any {
	events := p.pendingEvents
	p.pendingEvents = nil
	return events
}
