-- name: CreatePayment :exec
INSERT INTO payments (id, order_id, user_id, amount_cents, currency, status, payment_method, transaction_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetPaymentByID :one
SELECT id, order_id, user_id, amount_cents, currency, status, payment_method, transaction_id, created_at, updated_at
FROM payments
WHERE id = $1;

-- name: GetPaymentByOrderID :one
SELECT id, order_id, user_id, amount_cents, currency, status, payment_method, transaction_id, created_at, updated_at
FROM payments
WHERE order_id = $1;

-- name: UpdatePaymentStatus :exec
UPDATE payments
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