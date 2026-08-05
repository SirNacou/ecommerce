-- name: CreateProduct :exec
INSERT INTO products (id, category_id, name, description, price_cents, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetProductByID :one
SELECT id, category_id, name, description, price_cents, created_at, updated_at
FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT id, category_id, name, description, price_cents, created_at, updated_at
FROM products
WHERE (sqlc.narg(category_id)::uuid IS NULL OR category_id = sqlc.narg(category_id))
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListProductsByIds :many
SELECT id, category_id, name, description, price_cents, created_at, updated_at
FROM products
WHERE id = ANY($1::uuid[]);

-- name: CreateCategory :exec
INSERT INTO categories (id, name, slug, created_at)
VALUES ($1, $2, $3, $4);

-- name: GetCategoryByID :one
SELECT id, name, slug, created_at
FROM categories
WHERE id = $1;

-- name: ListCategories :many
SELECT id, name, slug, created_at
FROM categories
ORDER BY name ASC;

-- name: GetInventoryItem :one
SELECT product_id, available_quantity, reserved_quantity, created_at, updated_at
FROM inventory_items
WHERE product_id = $1;

-- name: GetInventoryItemForUpdate :one
SELECT product_id, available_quantity, reserved_quantity, created_at, updated_at
FROM inventory_items
WHERE product_id = $1
FOR UPDATE;

-- name: UpsertInventoryItem :exec
INSERT INTO inventory_items (product_id, available_quantity, reserved_quantity, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (product_id) DO UPDATE
SET available_quantity = EXCLUDED.available_quantity,
    reserved_quantity = EXCLUDED.reserved_quantity,
    updated_at = EXCLUDED.updated_at;

-- name: UpdateStockQuantities :exec
UPDATE inventory_items
SET available_quantity = $2,
    reserved_quantity = $3,
    updated_at = $4
WHERE product_id = $1;

-- name: CreateStockReservation :exec
INSERT INTO stock_reservations (id, product_id, quantity, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetStockReservation :one
SELECT id, product_id, quantity, status, created_at, updated_at
FROM stock_reservations
WHERE id = $1;

-- name: UpdateStockReservationStatus :exec
UPDATE stock_reservations
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