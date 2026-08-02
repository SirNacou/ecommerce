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