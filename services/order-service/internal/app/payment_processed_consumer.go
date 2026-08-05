package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SirNacou/ecommerce/services/order-service/internal/domain"
	"github.com/google/uuid"
)

// PaymentProcessedPayload is the payload of the PaymentProcessed event emitted
// by payment-service through the event bus.
type PaymentProcessedPayload struct {
	PaymentID     string `json:"payment_id"`
	OrderID       string `json:"order_id"`
	UserID        string `json:"user_id"`
	AmountCents   int64  `json:"amount_cents"`
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id"`
	Timestamp     string `json:"timestamp"`
}

type PaymentProcessedConsumer struct {
	uow UnitOfWork
}

func NewPaymentProcessedConsumer(uow UnitOfWork) *PaymentProcessedConsumer {
	return &PaymentProcessedConsumer{uow: uow}
}

func (c *PaymentProcessedConsumer) Handle(ctx context.Context, data []byte) error {
	var payload PaymentProcessedPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("parse PaymentProcessed payload: %w", err)
	}

	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		return fmt.Errorf("PaymentProcessed payload invalid order_id: %w", err)
	}

	var order *domain.Order
	err = c.uow.Execute(ctx, func(store OrderStore) error {
		var err error
		order, err = store.GetOrderByID(ctx, orderID)
		if err != nil {
			return err
		}

		if err := order.MarkAsPaid(payload.TransactionID); err != nil {
			return err
		}

		if err := store.UpdateOrderStatus(ctx, order.ID, order.Status); err != nil {
			return err
		}

		for _, event := range order.PopEvents() {
			eventPayload, err := json.Marshal(event)
			if err != nil {
				return err
			}
			if err := store.SaveOutboxEvent(ctx, "Order", order.ID.String(), event.EventType(), eventPayload); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
