package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/cart-service/internal/domain"
	"github.com/google/uuid"
)

type CartStore interface {
	GetCartByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error)
	CreateCart(ctx context.Context, cart *domain.Cart) error
	GetCartItems(ctx context.Context, cartID uuid.UUID) ([]*domain.CartItem, error)
	UpsertItem(ctx context.Context, cartID, productID uuid.UUID, quantity int32, priceCents int64) error
	UpdateItemQuantity(ctx context.Context, cartID, productID uuid.UUID, quantity int32) error
	RemoveItem(ctx context.Context, cartID, productID uuid.UUID) error
	ClearCart(ctx context.Context, cartID uuid.UUID) error
}

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(store CartStore) error) error
}
