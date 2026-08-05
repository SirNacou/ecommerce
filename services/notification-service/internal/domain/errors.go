package domain

import "errors"

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrInvalidRecipient     = errors.New("recipient address cannot be empty")
	ErrInvalidBody          = errors.New("notification body cannot be empty")
	ErrInvalidChannel       = errors.New("invalid or unsupported notification channel")
)
