-- +goose Up
-- Merge inventory-service into catalog-service on a single shared database.

-- Single source of truth for stock lives in inventory_items, so the
-- denormalized stock_quantity on products is removed.
ALTER TABLE products DROP COLUMN IF EXISTS stock_quantity;

CREATE TABLE inventory_items (
  product_id UUID PRIMARY KEY REFERENCES products (id) ON DELETE CASCADE,
  available_quantity INT NOT NULL DEFAULT 0 CHECK (available_quantity >= 0),
  reserved_quantity INT NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stock_reservations (
  id UUID PRIMARY KEY,
  product_id UUID NOT NULL REFERENCES inventory_items (product_id) ON DELETE CASCADE,
  quantity INT NOT NULL CHECK (quantity > 0),
  status VARCHAR(50) NOT NULL DEFAULT 'RESERVED', -- RESERVED, CONFIRMED, CANCELLED
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Unify the outbox schema with the other services (app-generated UUID PK,
-- PENDING/PROCESSED status). The old table is unused, so it can be recreated.
DROP TABLE IF EXISTS outbox_events;

CREATE TABLE outbox_events (
  id UUID PRIMARY KEY,
  aggregate_type VARCHAR(255) NOT NULL,
  aggregate_id VARCHAR(255) NOT NULL,
  event_type VARCHAR(255) NOT NULL,
  payload JSONB NOT NULL,
  status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_outbox_events_status ON outbox_events (status, created_at);

-- +goose Down
DROP TABLE IF EXISTS outbox_events;

DROP TABLE IF EXISTS stock_reservations;

DROP TABLE IF EXISTS inventory_items;

ALTER TABLE products ADD COLUMN stock_quantity INT NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0);