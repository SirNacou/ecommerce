package grpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/SirNacou/ecommerce/pkg/auth"
	v1 "github.com/SirNacou/ecommerce/services/payment-service/gen/v1"
	"github.com/SirNacou/ecommerce/services/payment-service/gen/v1/paymentv1connect"
	"github.com/SirNacou/ecommerce/services/payment-service/internal/app"
	"github.com/SirNacou/ecommerce/services/payment-service/internal/domain"
)

type PaymentHandler struct {
	paymentv1connect.UnimplementedPaymentServiceHandler
	processCmd *app.ProcessPaymentCommandHandler
	getQry     *app.GetPaymentQueryHandler
	refundCmd  *app.RefundPaymentCommandHandler
}

func NewPaymentHandler(
	processCmd *app.ProcessPaymentCommandHandler,
	getQry *app.GetPaymentQueryHandler,
	refundCmd *app.RefundPaymentCommandHandler,
) *PaymentHandler {
	return &PaymentHandler{
		processCmd: processCmd,
		getQry:     getQry,
		refundCmd:  refundCmd,
	}
}

func (h *PaymentHandler) ProcessPayment(
	ctx context.Context,
	req *connect.Request[v1.ProcessPaymentRequest],
) (*connect.Response[v1.ProcessPaymentResponse], error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	cmd := app.ProcessPaymentCommand{
		OrderID:       req.Msg.GetOrderId(),
		UserID:        userID,
		AmountCents:   req.Msg.GetAmountCents(),
		Currency:      req.Msg.GetCurrency(),
		PaymentMethod: req.Msg.GetPaymentMethod(),
	}

	payment, err := h.processCmd.Handle(ctx, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidAmount) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		if errors.Is(err, domain.ErrAlreadyProcessed) {
			return nil, connect.NewError(connect.CodeAlreadyExists, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.ProcessPaymentResponse{Payment: toProtoPayment(payment)}), nil
}

func (h *PaymentHandler) GetPayment(
	ctx context.Context,
	req *connect.Request[v1.GetPaymentRequest],
) (*connect.Response[v1.GetPaymentResponse], error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	payment, err := h.getQry.Handle(ctx, req.Msg.GetId(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.GetPaymentResponse{Payment: toProtoPayment(payment)}), nil
}

func (h *PaymentHandler) RefundPayment(
	ctx context.Context,
	req *connect.Request[v1.RefundPaymentRequest],
) (*connect.Response[v1.RefundPaymentResponse], error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	cmd := app.RefundPaymentCommand{
		PaymentID: req.Msg.GetPaymentId(),
		UserID:    userID,
		Reason:    req.Msg.GetReason(),
	}

	payment, err := h.refundCmd.Handle(ctx, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		if errors.Is(err, domain.ErrCannotRefund) {
			return nil, connect.NewError(connect.CodeFailedPrecondition, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.RefundPaymentResponse{Payment: toProtoPayment(payment)}), nil
}

func toProtoPayment(p *domain.Payment) *v1.Payment {
	return &v1.Payment{
		Id:            p.ID.String(),
		OrderId:       p.OrderID.String(),
		UserId:        p.UserID.String(),
		AmountCents:   p.AmountCents,
		Currency:      p.Currency,
		Status:        string(p.Status),
		PaymentMethod: p.PaymentMethod,
		TransactionId: p.TransactionID,
		CreatedAt:     timestamppb.New(p.CreatedAt),
	}
}
