-- name: CreateProduct :exec
INSERT INTO products (id, category_id, name, description, price_cents, stock_quantity, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetProductByID :one
SELECT id, category_id, name, description, price_cents, stock_quantity, created_at, updated_at
FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT id, category_id, name, description, price_cents, stock_quantity, created_at, updated_at
FROM products
WHERE (sqlc.narg(category_id)::uuid IS NULL OR category_id = sqlc.narg(category_id))
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

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