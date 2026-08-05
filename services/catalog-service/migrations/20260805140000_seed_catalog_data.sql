-- +goose Up
-- Seed sample catalog data for local development.

-- Categories
INSERT INTO categories (id, name, slug) VALUES
  ('11111111-1111-4111-8111-111111111111', 'Electronics', 'electronics'),
  ('22222222-2222-4222-8222-222222222222', 'Books', 'books'),
  ('33333333-3333-4333-8333-333333333333', 'Home & Kitchen', 'home-kitchen'),
  ('44444444-4444-4444-8444-444444444444', 'Sports', 'sports');

-- Products
INSERT INTO products (id, category_id, name, description, price_cents) VALUES
  ('aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa1', '11111111-1111-4111-8111-111111111111', 'Wireless Headphones', 'Over-ear noise cancelling headphones with 30h battery.', 7999),
  ('aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa2', '11111111-1111-4111-8111-111111111111', 'USB-C Charging Cable', '2m braided cable, fast charging supported.', 1299),
  ('aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa3', '11111111-1111-4111-8111-111111111111', 'Bluetooth Speaker', 'Portable waterproof speaker, 12h playtime.', 5999),
  ('aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa4', '11111111-1111-4111-8111-111111111111', 'Mechanical Keyboard', 'Hot-swappable 75% layout with RGB backlight.', 11999),
  ('bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbb1', '22222222-2222-4222-8222-222222222222', 'The Pragmatic Programmer', 'Classic guide to software craftsmanship.', 3499),
  ('bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbb2', '22222222-2222-4222-8222-222222222222', 'Clean Architecture', 'Robert C. Martin on software design principles.', 4299),
  ('cccccccc-cccc-4ccc-8ccc-ccccccccccc1', '33333333-3333-4333-8333-333333333333', 'Stainless Steel Water Bottle', 'Insulated 750ml bottle keeps drinks cold 24h.', 2499),
  ('cccccccc-cccc-4ccc-8ccc-ccccccccccc2', '33333333-3333-4333-8333-333333333333', 'Non-Stick Frying Pan', '28cm pan with scratch-resistant coating.', 3199),
  ('dddddddd-dddd-4ddd-8ddd-ddddddddddd1', '44444444-4444-4444-8444-444444444444', 'Yoga Mat', 'Eco-friendly 6mm mat with carrying strap.', 1999),
  ('dddddddd-dddd-4ddd-8ddd-ddddddddddd2', '44444444-4444-4444-8444-444444444444', 'Adjustable Dumbbells', 'Pair of dumbbells adjustable from 2.5 to 24kg.', 8999);

-- Inventory
INSERT INTO inventory_items (product_id, available_quantity, reserved_quantity) VALUES
  ('aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa1', 50, 0),
  ('aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa2', 200, 0),
  ('aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa3', 75, 0),
  ('aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa4', 40, 0),
  ('bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbb1', 100, 0),
  ('bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbb2', 80, 0),
  ('cccccccc-cccc-4ccc-8ccc-ccccccccccc1', 120, 0),
  ('cccccccc-cccc-4ccc-8ccc-ccccccccccc2', 60, 0),
  ('dddddddd-dddd-4ddd-8ddd-ddddddddddd1', 90, 0),
  ('dddddddd-dddd-4ddd-8ddd-ddddddddddd2', 30, 0);

-- +goose Down
DELETE FROM inventory_items WHERE product_id IN (
  'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa1',
  'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa2',
  'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa3',
  'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa4',
  'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbb1',
  'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbb2',
  'cccccccc-cccc-4ccc-8ccc-ccccccccccc1',
  'cccccccc-cccc-4ccc-8ccc-ccccccccccc2',
  'dddddddd-dddd-4ddd-8ddd-ddddddddddd1',
  'dddddddd-dddd-4ddd-8ddd-ddddddddddd2'
);

DELETE FROM products WHERE id IN (
  'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa1',
  'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa2',
  'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa3',
  'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa4',
  'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbb1',
  'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbb2',
  'cccccccc-cccc-4ccc-8ccc-ccccccccccc1',
  'cccccccc-cccc-4ccc-8ccc-ccccccccccc2',
  'dddddddd-dddd-4ddd-8ddd-ddddddddddd1',
  'dddddddd-dddd-4ddd-8ddd-ddddddddddd2'
);

DELETE FROM categories WHERE id IN (
  '11111111-1111-4111-8111-111111111111',
  '22222222-2222-4222-8222-222222222222',
  '33333333-3333-4333-8333-333333333333',
  '44444444-4444-4444-8444-444444444444'
);
