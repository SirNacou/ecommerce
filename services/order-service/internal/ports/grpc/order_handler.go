package grpc

import (
	"context"
	"errors"
	"strconv"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/SirNacou/ecommerce/pkg/auth"
	v1 "github.com/SirNacou/ecommerce/services/order-service/gen/v1"
	"github.com/SirNacou/ecommerce/services/order-service/gen/v1/orderv1connect"
	"github.com/SirNacou/ecommerce/services/order-service/internal/app"
	"github.com/SirNacou/ecommerce/services/order-service/internal/domain"
)

type OrderHandler struct {
	orderv1connect.UnimplementedOrderServiceHandler
	createOrderCmd *app.CreateOrderCommandHandler
	getOrderQry    *app.GetOrderQueryHandler
	listOrdersQry  *app.ListOrdersQueryHandler
	cancelOrderCmd *app.CancelOrderCommandHandler
}

func NewOrderHandler(
	createOrderCmd *app.CreateOrderCommandHandler,
	getOrderQry *app.GetOrderQueryHandler,
	listOrdersQry *app.ListOrdersQueryHandler,
	cancelOrderCmd *app.CancelOrderCommandHandler,
) *OrderHandler {
	return &OrderHandler{
		createOrderCmd: createOrderCmd,
		getOrderQry:    getOrderQry,
		listOrdersQry:  listOrdersQry,
		cancelOrderCmd: cancelOrderCmd,
	}
}

func (h *OrderHandler) CreateOrder(
	ctx context.Context,
	req *connect.Request[v1.CreateOrderRequest],
) (*connect.Response[v1.CreateOrderResponse], error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	items := make([]app.OrderItemInput, 0, len(req.Msg.GetItems()))
	for _, item := range req.Msg.GetItems() {
		items = append(items, app.OrderItemInput{
			ProductID: item.GetProductId(),
			Quantity:  item.GetQuantity(),
		})
	}

	order, err := h.createOrderCmd.Handle(ctx, app.CreateOrderCommand{
		UserID: userID,
		Items:  items,
	})
	if err != nil {
		if errors.Is(err, domain.ErrEmptyOrder) || errors.Is(err, domain.ErrInvalidQuantity) || errors.Is(err, domain.ErrInvalidPrice) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.CreateOrderResponse{Order: toProtoOrder(order)}), nil
}

func (h *OrderHandler) GetOrder(
	ctx context.Context,
	req *connect.Request[v1.GetOrderRequest],
) (*connect.Response[v1.GetOrderResponse], error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	order, err := h.getOrderQry.Handle(ctx, req.Msg.GetId(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.GetOrderResponse{Order: toProtoOrder(order)}), nil
}

func (h *OrderHandler) ListOrders(
	ctx context.Context,
	req *connect.Request[v1.ListOrdersRequest],
) (*connect.Response[v1.ListOrdersResponse], error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	pageSize := req.Msg.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10
	}

	var offset int32 = 0
	if token := req.Msg.GetPageToken(); token != "" {
		if parsed, err := strconv.Atoi(token); err == nil && parsed >= 0 {
			offset = int32(parsed)
		}
	}

	orders, err := h.listOrdersQry.Handle(ctx, app.ListOrdersQuery{
		UserID:   userID,
		PageSize: pageSize,
		Offset:   offset,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoOrders := make([]*v1.Order, 0, len(orders))
	for _, o := range orders {
		protoOrders = append(protoOrders, toProtoOrder(o))
	}

	var nextPageToken string
	if int32(len(orders)) == pageSize {
		nextPageToken = strconv.Itoa(int(offset + pageSize))
	}

	return connect.NewResponse(&v1.ListOrdersResponse{
		Orders:        protoOrders,
		NextPageToken: nextPageToken,
	}), nil
}

func (h *OrderHandler) CancelOrder(
	ctx context.Context,
	req *connect.Request[v1.CancelOrderRequest],
) (*connect.Response[v1.CancelOrderResponse], error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	cmd := app.CancelOrderCommand{
		OrderID: req.Msg.GetId(),
		UserID:  userID,
		Reason:  req.Msg.GetReason(),
	}

	order, err := h.cancelOrderCmd.Handle(ctx, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		if errors.Is(err, domain.ErrCannotCancelOrder) {
			return nil, connect.NewError(connect.CodeFailedPrecondition, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.CancelOrderResponse{Order: toProtoOrder(order)}), nil
}

func toProtoOrder(o *domain.Order) *v1.Order {
	items := make([]*v1.OrderItem, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, &v1.OrderItem{
			Id:         item.ID.String(),
			ProductId:  item.ProductID.String(),
			Quantity:   item.Quantity,
			PriceCents: item.PriceCents,
		})
	}

	return &v1.Order{
		Id:         o.ID.String(),
		UserId:     o.UserID.String(),
		Items:      items,
		TotalCents: o.TotalCents,
		Status:     string(o.Status),
		CreatedAt:  timestamppb.New(o.CreatedAt),
	}
}
