-- +goose Up
CREATE TABLE carts (
  id UUID PRIMARY KEY,
  user_id UUID UNIQUE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cart_items (
  id UUID PRIMARY KEY,
  cart_id UUID NOT NULL REFERENCES carts (id) ON DELETE CASCADE,
  product_id UUID NOT NULL,
  quantity INT NOT NULL CHECK (quantity > 0),
  price_cents BIGINT NOT NULL CHECK (price_cents >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT unique_cart_product UNIQUE (cart_id, product_id)
);

CREATE INDEX idx_cart_user_id ON carts (user_id);

CREATE INDEX idx_cart_items_cart_id ON cart_items (cart_id);

-- +goose Down
DROP TABLE IF EXISTS cart_items;

DROP TABLE IF EXISTS carts;