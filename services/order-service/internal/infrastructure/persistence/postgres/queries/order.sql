-- name: CreateOrder :exec
INSERT INTO orders (id, user_id, total_cents, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: CreateOrderItem :exec
INSERT INTO order_items (id, order_id, product_id, quantity, price_cents, created_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetOrderByID :one
SELECT id, user_id, total_cents, status, created_at, updated_at
FROM orders
WHERE id = $1;

-- name: GetOrderItemsByOrderID :many
SELECT id, order_id, product_id, quantity, price_cents, created_at
FROM order_items
WHERE order_id = $1;

-- name: ListOrdersByUserID :many
SELECT id, user_id, total_cents, status, created_at, updated_at
FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateOrderStatus :exec
UPDATE orders
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