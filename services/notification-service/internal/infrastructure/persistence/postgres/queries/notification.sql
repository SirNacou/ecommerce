-- name: CreateNotification :exec
INSERT INTO notifications (id, user_id, channel, recipient, subject, body, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetNotificationByID :one
SELECT id, user_id, channel, recipient, subject, body, status, created_at, updated_at
FROM notifications
WHERE id = $1;

-- name: ListNotificationsByUserID :many
SELECT id, user_id, channel, recipient, subject, body, status, created_at, updated_at
FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateNotificationStatus :exec
UPDATE notifications
SET status = $2,
    updated_at = $3
WHERE id = $1;

-- name: InsertOutboxEvent :exec
INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload, status, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetPendingOutboxEvents :many
SELECT id, aggregate_type, aggregate_id, event_type, payload, status, created_at
FROM outbox_events
WHERE status = 'PENDING'
ORDER BY created_at ASC
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: MarkOutboxEventProcessed :exec
UPDATE outbox_events
SET status = 'PROCESSED'
WHERE id = $1;