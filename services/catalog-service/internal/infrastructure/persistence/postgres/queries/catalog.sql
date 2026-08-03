-- name: CreateCategory :one
INSERT INTO categories (id, name, slug)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListCategories :many
SELECT * FROM categories ORDER BY name ASC;

-- name: CreateProduct :one
INSERT INTO products (id, category_id, name, description, price_cents, stock_quantity)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: ListProducts :many
SELECT * FROM products
WHERE (sqlc.narg('category_id')::uuid IS NULL OR category_id = sqlc.narg('category_id'))
ORDER BY created_at DESC
LIMIT $1;

-- name: CreateOutboxEvent :one
INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;