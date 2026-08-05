package domain

import (
	"time"

	"github.com/google/uuid"
)

type Channel string

const (
	ChannelEmail Channel = "EMAIL"
	ChannelSMS   Channel = "SMS"
	ChannelPush  Channel = "PUSH"
)

type NotificationStatus string

const (
	StatusPending NotificationStatus = "PENDING"
	StatusSent    NotificationStatus = "SENT"
	StatusFailed  NotificationStatus = "FAILED"
)

type Notification struct {
	AggregateRoot
	ID        uuid.UUID
	UserID    uuid.UUID
	Channel   Channel
	Recipient string
	Subject   string
	Body      string
	Status    NotificationStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NotificationSentEvent struct {
	NotificationID string    `json:"notification_id"`
	UserID         string    `json:"user_id"`
	Channel        string    `json:"channel"`
	Recipient      string    `json:"recipient"`
	Timestamp      time.Time `json:"timestamp"`
}

func (e NotificationSentEvent) EventType() string     { return "NotificationSent" }
func (e NotificationSentEvent) OccurredAt() time.Time { return e.Timestamp }

func NewNotification(
	userID uuid.UUID,
	channel Channel,
	recipient, subject, body string,
) (*Notification, error) {
	if recipient == "" {
		return nil, ErrInvalidRecipient
	}
	if body == "" {
		return nil, ErrInvalidBody
	}
	if channel != ChannelEmail && channel != ChannelSMS && channel != ChannelPush {
		return nil, ErrInvalidChannel
	}

	id := uuid.New()
	now := time.Now().UTC()

	notification := &Notification{
		ID:        id,
		UserID:    userID,
		Channel:   channel,
		Recipient: recipient,
		Subject:   subject,
		Body:      body,
		Status:    StatusSent,
		CreatedAt: now,
		UpdatedAt: now,
	}

	notification.RecordEvent(NotificationSentEvent{
		NotificationID: id.String(),
		UserID:         userID.String(),
		Channel:        string(channel),
		Recipient:      recipient,
		Timestamp:      now,
	})

	return notification, nil
}
