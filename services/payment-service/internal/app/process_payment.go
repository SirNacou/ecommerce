package app

import (
	"context"
	"encoding/json"

	"github.com/SirNacou/ecommerce/services/payment-service/internal/domain"
	"github.com/google/uuid"
)

type ProcessPaymentCommand struct {
	OrderID       string
	UserID        string
	AmountCents   int64
	Currency      string
	PaymentMethod string
}

type ProcessPaymentCommandHandler struct {
	uow UnitOfWork
}

func NewProcessPaymentCommandHandler(uow UnitOfWork) *ProcessPaymentCommandHandler {
	return &ProcessPaymentCommandHandler{uow: uow}
}

func (h *ProcessPaymentCommandHandler) Handle(ctx context.Context, cmd ProcessPaymentCommand) (*domain.Payment, error) {
	orderID, err := uuid.Parse(cmd.OrderID)
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}

	payment, err := domain.NewPayment(orderID, userID, cmd.AmountCents, cmd.Currency, cmd.PaymentMethod)
	if err != nil {
		return nil, err
	}

	err = h.uow.Execute(ctx, func(store PaymentStore) error {
		existing, err := store.GetPaymentByOrderID(ctx, orderID)
		if err == nil && existing != nil {
			return domain.ErrAlreadyProcessed
		}

		if err := store.CreatePayment(ctx, payment); err != nil {
			return err
		}

		for _, event := range payment.PopEvents() {
			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}
			if err := store.SaveOutboxEvent(ctx, "Payment", payment.ID.String(), event.EventType(), payload); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return payment, nil
}
