ALTER TABLE users DROP COLUMN IF EXISTS product_sold_total;
ALTER TABLE users ADD COLUMN IF NOT EXISTS product_sold_total INT NOT NULL DEFAULT 0;
