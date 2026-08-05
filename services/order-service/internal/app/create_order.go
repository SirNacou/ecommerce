package app

import (
	"context"
	"encoding/json"

	"github.com/SirNacou/ecommerce/services/order-service/internal/domain"
	"github.com/google/uuid"
)

type OrderItemInput struct {
	ProductID string
	Quantity  int32
}

type CreateOrderCommand struct {
	UserID string
	Items  []OrderItemInput
}

// PriceResolver resolves authoritative product prices for a set of product ids.
type PriceResolver interface {
	ResolvePrices(ctx context.Context, productIDs []string) (map[string]int64, error)
}

type CreateOrderCommandHandler struct {
	uow    UnitOfWork
	prices PriceResolver
}

func NewCreateOrderCommandHandler(uow UnitOfWork, prices PriceResolver) *CreateOrderCommandHandler {
	return &CreateOrderCommandHandler{uow: uow, prices: prices}
}

func (h *CreateOrderCommandHandler) Handle(ctx context.Context, cmd CreateOrderCommand) (*domain.Order, error) {
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}

	productIDs := make([]string, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		if _, err := uuid.Parse(item.ProductID); err != nil {
			return nil, domain.ErrInvalidQuantity
		}
		productIDs = append(productIDs, item.ProductID)
	}

	prices, err := h.prices.ResolvePrices(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	inputs := make([]struct {
		ProductID  uuid.UUID
		Quantity   int32
		PriceCents int64
	}, 0, len(cmd.Items))

	for _, item := range cmd.Items {
		priceCents, ok := prices[item.ProductID]
		if !ok {
			return nil, domain.ErrProductNotFound
		}
		prodID, _ := uuid.Parse(item.ProductID)
		inputs = append(inputs, struct {
			ProductID  uuid.UUID
			Quantity   int32
			PriceCents int64
		}{
			ProductID:  prodID,
			Quantity:   item.Quantity,
			PriceCents: priceCents,
		})
	}

	order, err := domain.NewOrder(userID, inputs)
	if err != nil {
		return nil, err
	}

	err = h.uow.Execute(ctx, func(store OrderStore) error {
		if err := store.CreateOrder(ctx, order); err != nil {
			return err
		}

		for _, event := range order.PopEvents() {
			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}
			if err := store.SaveOutboxEvent(ctx, "Order", order.ID.String(), event.EventType(), payload); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}