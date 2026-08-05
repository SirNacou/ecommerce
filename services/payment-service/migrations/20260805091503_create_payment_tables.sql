-- +goose Up
CREATE TABLE payments (
  id UUID PRIMARY KEY,
  order_id UUID NOT NULL UNIQUE,
  user_id UUID NOT NULL,
  amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
  currency VARCHAR(3) NOT NULL DEFAULT 'USD',
  status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
  payment_method VARCHAR(50) NOT NULL,
  transaction_id VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE outbox_events (
  id UUID PRIMARY KEY,
  aggregate_type VARCHAR(255) NOT NULL,
  aggregate_id VARCHAR(255) NOT NULL,
  event_type VARCHAR(255) NOT NULL,
  payload JSONB NOT NULL,
  status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payments_order_id ON payments (order_id);

CREATE INDEX idx_payments_user_id ON payments (user_id);

CREATE INDEX idx_outbox_events_status ON outbox_events (status, created_at);

-- +goose Down
DROP TABLE IF EXISTS outbox_events;

DROP TABLE IF EXISTS payments;