package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SirNacou/ecommerce/services/notification-service/internal/domain"
)

// OrderCreatedPayload is the payload of the OrderCreated event emitted by
// order-service through the event bus.
type OrderCreatedPayload struct {
	OrderID    string `json:"order_id"`
	UserID     string `json:"user_id"`
	TotalCents int64  `json:"total_cents"`
	Timestamp  string `json:"timestamp"`
}

type OrderCreatedConsumer struct {
	sendCmd *SendNotificationCommandHandler
}

func NewOrderCreatedConsumer(sendCmd *SendNotificationCommandHandler) *OrderCreatedConsumer {
	return &OrderCreatedConsumer{sendCmd: sendCmd}
}

func (c *OrderCreatedConsumer) Handle(ctx context.Context, data []byte) error {
	var payload OrderCreatedPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("parse OrderCreated payload: %w", err)
	}

	if payload.UserID == "" {
		return fmt.Errorf("OrderCreated payload missing user_id")
	}

	_, err := c.sendCmd.Handle(ctx, SendNotificationCommand{
		UserID:    payload.UserID,
		Channel:   domain.ChannelEmail,
		Recipient: payload.UserID,
		Subject:   "Order Confirmation",
		Body:      fmt.Sprintf("Your order %s has been placed (total %d cents).", payload.OrderID, payload.TotalCents),
	})
	if err != nil {
		return fmt.Errorf("send order notification: %w", err)
	}

	return nil
}