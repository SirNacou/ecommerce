package app

import (
	"context"
	"encoding/json"

	"github.com/SirNacou/ecommerce/services/payment-service/internal/domain"
	"github.com/google/uuid"
)

type RefundPaymentCommand struct {
	PaymentID string
	UserID    string
	Reason    string
}

type RefundPaymentCommandHandler struct {
	uow UnitOfWork
}

func NewRefundPaymentCommandHandler(uow UnitOfWork) *RefundPaymentCommandHandler {
	return &RefundPaymentCommandHandler{uow: uow}
}

func (h *RefundPaymentCommandHandler) Handle(ctx context.Context, cmd RefundPaymentCommand) (*domain.Payment, error) {
	paymentID, err := uuid.Parse(cmd.PaymentID)
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}

	var payment *domain.Payment
	err = h.uow.Execute(ctx, func(store PaymentStore) error {
		var err error
		payment, err = store.GetPaymentByID(ctx, paymentID)
		if err != nil {
			return err
		}

		if payment.UserID != userID {
			return domain.ErrPaymentNotFound
		}

		if err := payment.Refund(cmd.Reason); err != nil {
			return err
		}

		if err := store.UpdatePaymentStatus(ctx, payment.ID, payment.Status); err != nil {
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
