package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type ReserveStockCommand struct {
	ReservationID string
	ProductID     string
	Quantity      int32
}

type ReserveStockCommandHandler struct {
	uow UnitOfWork
}

func NewReserveStockCommandHandler(uow UnitOfWork) *ReserveStockCommandHandler {
	return &ReserveStockCommandHandler{uow: uow}
}

func (h *ReserveStockCommandHandler) Handle(ctx context.Context, cmd ReserveStockCommand) error {
	resID, err := uuid.Parse(cmd.ReservationID)
	if err != nil {
		return fmt.Errorf("invalid reservation id: %w", err)
	}

	prodID, err := uuid.Parse(cmd.ProductID)
	if err != nil {
		return fmt.Errorf("invalid product id: %w", err)
	}

	return h.uow.Execute(ctx, func(store CatalogStore) error {
		item, err := store.GetItemForUpdate(ctx, prodID)
		if err != nil {
			return err
		}

		if err := item.Reserve(cmd.Quantity); err != nil {
			return err
		}

		if err := store.UpdateStock(ctx, item); err != nil {
			return err
		}

		if err := store.CreateReservation(ctx, resID, prodID, cmd.Quantity); err != nil {
			return err
		}

		payload, _ := json.Marshal(map[string]any{
			"reservation_id": resID.String(),
			"product_id":     prodID.String(),
			"quantity":       cmd.Quantity,
		})

		return store.SaveOutboxEvent(ctx, "InventoryItem", prodID.String(), "StockReserved", payload)
	})
}
