package app

import (
	"context"

	"github.com/SirNacou/ecommerce/services/payment-service/internal/domain"
	"github.com/google/uuid"
)

type GetPaymentQueryHandler struct {
	uow UnitOfWork
}

func NewGetPaymentQueryHandler(uow UnitOfWork) *GetPaymentQueryHandler {
	return &GetPaymentQueryHandler{uow: uow}
}

func (h *GetPaymentQueryHandler) Handle(ctx context.Context, idStr, userIDStr string) (*domain.Payment, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}

	var payment *domain.Payment
	err = h.uow.Execute(ctx, func(store PaymentStore) error {
		var err error
		payment, err = store.GetPaymentByID(ctx, id)
		return err
	})
	if err != nil {
		return nil, err
	}

	if payment.UserID != userID {
		return nil, domain.ErrPaymentNotFound
	}

	return payment, nil
}
