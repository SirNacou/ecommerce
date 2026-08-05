package grpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/SirNacou/ecommerce/pkg/auth"
	v1 "github.com/SirNacou/ecommerce/services/notification-service/gen/v1"
	"github.com/SirNacou/ecommerce/services/notification-service/gen/v1/notificationv1connect"
	"github.com/SirNacou/ecommerce/services/notification-service/internal/app"
	"github.com/SirNacou/ecommerce/services/notification-service/internal/domain"
)

type NotificationHandler struct {
	notificationv1connect.UnimplementedNotificationServiceHandler
	sendCmd *app.SendNotificationCommandHandler
	getQry  *app.GetNotificationQueryHandler
}

func NewNotificationHandler(
	sendCmd *app.SendNotificationCommandHandler,
	getQry *app.GetNotificationQueryHandler,
) *NotificationHandler {
	return &NotificationHandler{
		sendCmd: sendCmd,
		getQry:  getQry,
	}
}

func (h *NotificationHandler) SendNotification(
	ctx context.Context,
	req *connect.Request[v1.SendNotificationRequest],
) (*connect.Response[v1.SendNotificationResponse], error) {
	_, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	cmd := app.SendNotificationCommand{
		UserID:    req.Msg.GetUserId(),
		Channel:   mapProtoChannel(req.Msg.GetChannel()),
		Recipient: req.Msg.GetRecipient(),
		Subject:   req.Msg.GetSubject(),
		Body:      req.Msg.GetBody(),
	}

	notification, err := h.sendCmd.Handle(ctx, cmd)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidRecipient) || errors.Is(err, domain.ErrInvalidBody) || errors.Is(err, domain.ErrInvalidChannel) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.SendNotificationResponse{Notification: toProtoNotification(notification)}), nil
}

func (h *NotificationHandler) GetNotification(
	ctx context.Context,
	req *connect.Request[v1.GetNotificationRequest],
) (*connect.Response[v1.GetNotificationResponse], error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	notification, err := h.getQry.Handle(ctx, req.Msg.GetId(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotificationNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.GetNotificationResponse{Notification: toProtoNotification(notification)}), nil
}

func mapProtoChannel(c v1.Channel) domain.Channel {
	switch c {
	case v1.Channel_CHANNEL_EMAIL:
		return domain.ChannelEmail
	case v1.Channel_CHANNEL_SMS:
		return domain.ChannelSMS
	case v1.Channel_CHANNEL_PUSH:
		return domain.ChannelPush
	default:
		return ""
	}
}

func toProtoNotification(n *domain.Notification) *v1.Notification {
	var ch v1.Channel
	switch n.Channel {
	case domain.ChannelEmail:
		ch = v1.Channel_CHANNEL_EMAIL
	case domain.ChannelSMS:
		ch = v1.Channel_CHANNEL_SMS
	case domain.ChannelPush:
		ch = v1.Channel_CHANNEL_PUSH
	}

	return &v1.Notification{
		Id:        n.ID.String(),
		UserId:    n.UserID.String(),
		Channel:   ch,
		Recipient: n.Recipient,
		Subject:   n.Subject,
		Body:      n.Body,
		Status:    string(n.Status),
		CreatedAt: timestamppb.New(n.CreatedAt),
	}
}