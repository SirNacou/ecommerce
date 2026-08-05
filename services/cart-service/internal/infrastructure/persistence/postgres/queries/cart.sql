-- name: GetCartByUserID :one
SELECT id, user_id, created_at, updated_at
FROM carts
WHERE user_id = $1;

-- name: CreateCart :exec
INSERT INTO carts (id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4);

-- name: GetCartItems :many
SELECT id, cart_id, product_id, quantity, price_cents, created_at, updated_at
FROM cart_items
WHERE cart_id = $1
ORDER BY created_at ASC;

-- name: UpsertCartItem :exec
INSERT INTO cart_items (id, cart_id, product_id, quantity, price_cents, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (cart_id, product_id) DO UPDATE
SET quantity = cart_items.quantity + EXCLUDED.quantity,
    price_cents = EXCLUDED.price_cents,
    updated_at = EXCLUDED.updated_at;

-- name: UpdateCartItemQuantity :exec
UPDATE cart_items
SET quantity = $3,
    updated_at = $4
WHERE cart_id = $1 AND product_id = $2;

-- name: RemoveCartItem :exec
DELETE FROM cart_items
WHERE cart_id = $1 AND product_id = $2;

-- name: ClearCartItems :exec
DELETE FROM cart_items
WHERE cart_id = $1;