package grpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/SirNacou/ecommerce/pkg/auth"
	v1 "github.com/SirNacou/ecommerce/services/cart-service/gen/v1"
	"github.com/SirNacou/ecommerce/services/cart-service/gen/v1/cartv1connect"
	"github.com/SirNacou/ecommerce/services/cart-service/internal/app"
	"github.com/SirNacou/ecommerce/services/cart-service/internal/domain"
)

type CartHandler struct {
	cartv1connect.UnimplementedCartServiceHandler
	getCartQry    *app.GetCartQueryHandler
	addItemCmd    *app.AddItemCommandHandler
	updateQtyCmd  *app.UpdateQuantityCommandHandler
	removeItemCmd *app.RemoveItemCommandHandler
	clearCartCmd  *app.ClearCartCommandHandler
}

func NewCartHandler(
	getCartQry *app.GetCartQueryHandler,
	addItemCmd *app.AddItemCommandHandler,
	updateQtyCmd *app.UpdateQuantityCommandHandler,
	removeItemCmd *app.RemoveItemCommandHandler,
	clearCartCmd *app.ClearCartCommandHandler,
) *CartHandler {
	return &CartHandler{
		getCartQry:    getCartQry,
		addItemCmd:    addItemCmd,
		updateQtyCmd:  updateQtyCmd,
		removeItemCmd: removeItemCmd,
		clearCartCmd:  clearCartCmd,
	}
}

func (h *CartHandler) GetCart(
	ctx context.Context,
	req *connect.Request[v1.GetCartRequest],
) (*connect.Response[v1.GetCartResponse], error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	cart, err := h.getCartQry.Handle(ctx, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.GetCartResponse{Cart: toProtoCart(cart)}), nil
}

func (h *CartHandler) AddItem(
	ctx context.Context,
	req *connect.Request[v1.AddItemRequest],
) (*connect.Response[v1.AddItemResponse], error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	cmd := app.AddItemCommand{
		UserID:    userID,
		ProductID: req.Msg.GetProductId(),
		Quantity:  req.Msg.GetQuantity(),
	}

	cart, err := h.addItemCmd.Handle(ctx, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidQuantity) || errors.Is(err, domain.ErrProductNotFound) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.AddItemResponse{Cart: toProtoCart(cart)}), nil
}

func (h *CartHandler) UpdateItemQuantity(
	ctx context.Context,
	req *connect.Request[v1.UpdateItemQuantityRequest],
) (*connect.Response[v1.UpdateItemQuantityResponse], error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	cmd := app.UpdateQuantityCommand{
		UserID:    userID,
		ProductID: req.Msg.GetProductId(),
		Quantity:  req.Msg.GetQuantity(),
	}

	cart, err := h.updateQtyCmd.Handle(ctx, cmd)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.UpdateItemQuantityResponse{Cart: toProtoCart(cart)}), nil
}

func (h *CartHandler) RemoveItem(
	ctx context.Context,
	req *connect.Request[v1.RemoveItemRequest],
) (*connect.Response[v1.RemoveItemResponse], error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	cart, err := h.removeItemCmd.Handle(ctx, userID, req.Msg.GetProductId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.RemoveItemResponse{Cart: toProtoCart(cart)}), nil
}

func (h *CartHandler) ClearCart(
	ctx context.Context,
	req *connect.Request[v1.ClearCartRequest],
) (*connect.Response[v1.ClearCartResponse], error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	if err := h.clearCartCmd.Handle(ctx, userID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.ClearCartResponse{Success: true}), nil
}

func toProtoCart(c *domain.Cart) *v1.Cart {
	items := make([]*v1.CartItem, 0, len(c.Items))
	for _, item := range c.Items {
		items = append(items, &v1.CartItem{
			Id:         item.ID.String(),
			ProductId:  item.ProductID.String(),
			Quantity:   item.Quantity,
			PriceCents: item.PriceCents,
		})
	}

	return &v1.Cart{
		Id:         c.ID.String(),
		UserId:     c.UserID.String(),
		Items:      items,
		TotalCents: c.TotalCents(),
	}
}
