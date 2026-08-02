-- +goose Up
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users (email);

-- +goose Down
DROP TABLE IF EXISTS users;