-- name: CreateUser :exec
INSERT INTO users (id, email, password_hash, name, created_at)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateUser :exec
UPDATE users
SET email = $2, password_hash = $3, name = $4
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, name, created_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, password_hash, name, created_at
FROM users
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